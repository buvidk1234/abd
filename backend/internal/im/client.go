package im

import (
	"backend/internal/api/apiresp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	mu   sync.Mutex
	conn *websocket.Conn

	respWriter http.ResponseWriter
	req        *http.Request

	PlatformID int
	UserID     string
	IsCompress bool
	Encoder

	server *WsServer

	token     string
	closed    atomic.Bool
	closedErr error

	hbCtx    context.Context
	hbCancel context.CancelFunc
}

// typed context keys to avoid collisions
type ctxKey string

const (
	ctxKeySendID   ctxKey = "send_id"
	ctxKeyPlatform ctxKey = "platform_id"
)

func (c *Client) ResetClient(respWriter http.ResponseWriter, req *http.Request, conn *websocket.Conn, wsServer *WsServer) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.respWriter = respWriter
	c.req = req
	c.conn = conn
	c.server = wsServer
	// parse URL parameters
	c.UserID = req.URL.Query().Get(WsUserID)
	c.PlatformID, _ = strconv.Atoi(req.URL.Query().Get(PlatformID))
	c.IsCompress = req.URL.Query().Get(Compression) == GzipCompressionProtocol
	c.token = req.URL.Query().Get(Token)

	c.closed.Store(false)
	c.hbCtx, c.hbCancel = context.WithCancel(req.Context())

	c.Encoder = NewJsonEncoder()
}

func (c *Client) PushUserOnlineStatus(data []byte) error {
	return c.writeBinaryMsg(Resp{
		ReqIdentifier: WsSubUserOnlineStatus,
		Data:          data,
	})
}

func (c *Client) KickOnlineMessage() error {
	resp := Resp{
		ReqIdentifier: WSKickOnlineMsg,
	}
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteMessage(PongMessage, []byte(appData))
}

func (c *Client) pongHandler(_ string) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return err
	}
	return nil
}

func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("socket have panic err: %v\n%s", r, debug.Stack())
			c.closedErr = errors.New("panic err")
		}
		c.close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	c.activeHeartbeat(c.hbCtx)

	for {
		messageType, message, err := c.conn.ReadMessage()

		if c.closed.Load() {
			log.Printf("connection is closed: %v", c.req.Context())
			c.closedErr = errors.New("connection is closed")
			return
		}

		if err != nil {
			log.Printf("readMessage: %v", err)
			return
		}
		log.Println("读取到信息")

		switch messageType {
		case MessageText:
			_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
			c.handleMessage(message)
		//case MessageText:
		//	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		//	if err := c.handleTextMessage(message); err != nil {
		//		log.Printf("handleTextMessage: %v", err)
		//		return
		//	}
		case PingMessage:
			c.conn.WriteMessage(PongMessage, nil)
		case CloseMessage:
			c.closedErr = errors.New("client actively close the connection")
			return
		default:
		}
	}
}

func (c *Client) handleMessage(b []byte) error {
	if c.IsCompress {
		var err error
		b, err = c.server.Decompress(b)
		if err != nil {
			return err
		}
	}

	binaryReq := getReq()
	defer freeReq(binaryReq)

	if err := c.Encoder.Decode(b, binaryReq); err != nil {
		return err
	}

	var (
		resp any
		err  error
	)

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeySendID, binaryReq.SendID)
	ctx = context.WithValue(ctx, ctxKeyPlatform, c.PlatformID)
	log.Print("调用后端服务")
	switch binaryReq.ReqIdentifier {
	case WSSendMsg:
		log.Println("发信息")
		resp, err = c.server.SendMessage(ctx, binaryReq)
	case WSPullSpecifiedConv:
		log.Println("拉取指定会话")
		resp, err = c.server.PullSpecifiedConv(ctx, binaryReq)
	case WSPullConvList:
		log.Println("拉取会话列表")
		resp, err = c.server.PullConvList(ctx, binaryReq)
	default:
		log.Printf("unknown req identifier: %d", binaryReq.ReqIdentifier)
		return fmt.Errorf(
			"ReqIdentifier failed,sendID:%s,msgIncr:%s,reqIdentifier:%d",
			binaryReq.SendID,
			binaryReq.MsgIncr,
			binaryReq.ReqIdentifier,
		)

	}
	return c.replyMessage(binaryReq, resp, err)
}

func (c *Client) replyMessage(req *Req, resp any, err error) error {
	errResp := apiresp.ParseError(err)
	reply := Resp{
		ReqIdentifier: req.ReqIdentifier,
		MsgIncr:       req.MsgIncr,
		Code:          errResp.Code,
		Msg:           errResp.Msg,
		Data:          resp,
	}
	return c.writeBinaryMsg(reply)
}

func (c *Client) writeBinaryMsg(resp Resp) error {
	data, err := c.Encoder.Encode(resp)
	if err != nil {
		return err
	}
	//if c.IsCompress {
	//	data, err = c.server.Compress(data)
	//	if err != nil {
	//		return err
	//	}
	//}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err = c.conn.WriteMessage(MessageText, data)
	return err
}

func (c *Client) handleTextMessage(b []byte) error {
	var msg TextMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return err
	}
	switch msg.Type {
	case TextPong:
		return nil
	case TextPing:
		msg.Type = TextPong
		respData, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return c.safeWriteMessage(MessageText, respData)
	default:
		return fmt.Errorf("not support message type %s", msg.Type)
	}
}

func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	_ = c.conn.Close()
	c.hbCancel()
	c.server.UnRegister(c)
}

func (c *Client) PushMessage(ctx context.Context, msg any) error {
	resp := Resp{
		ReqIdentifier: WSPushMsg,
		Data:          msg,
	}
	return c.writeBinaryMsg(resp)
}

func (c *Client) activeHeartbeat(ctx context.Context) {
	if c.PlatformID == WebPlatformID {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("activeHeartbeat panic user=%s platform=%d: %v\n%s", c.UserID, c.PlatformID, r, debug.Stack())
				}
			}()

			log.Printf("server initiative send heartbeat start. user=%s platform=%d", c.UserID, c.PlatformID)

			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					if err := c.safeWriteMessage(PingMessage, nil); err != nil {
						log.Printf("send Ping Message error: %v", err)
						return
					}
				case <-c.hbCtx.Done():
					return
				}
			}
		}()
	}
}

func (c *Client) safeWriteMessage(messageType int, data []byte) error {
	if c.closed.Load() {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, data)
}
