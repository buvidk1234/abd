package im

import (
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/pkg/prommetrics"
	"backend/internal/service"
	"backend/pkg/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

// Config WebSocket 服务配置
type Config struct {
	Addr             string `yaml:"addr"`              // WebSocket 监听地址
	MaxConnNum       int64  `yaml:"max_conn_num"`      // 最大连接数
	WriteBufferSize  int    `yaml:"write_buffer_size"` // 写缓冲区大小
	HandshakeTimeout int    `yaml:"handshake_timeout"` // 握手超时(秒)
}

type WsServer struct {
	addr              string
	wsMaxConnNum      int64
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	Clients           UserMap
	clientPool        sync.Pool
	handshakeTimeout  time.Duration
	writeBufferSize   int
	// ready             atomic.Bool

	registerChan    chan *Client
	unregisterChan  chan *Client
	kickHandlerChan chan *kickHandler
	validate        *validator.Validate
	Compressor
	MessageHandler

	authClient *service.UserService
}

type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func NewWsServer(cfg Config) *WsServer {
	producer, err := kafka.NewSyncProducer()
	if err != nil {
		panic(fmt.Sprintf("failed to create kafka producer: %v", err))
	}

	// 设置默认值
	if cfg.Addr == "" {
		cfg.Addr = ":8082"
	}
	if cfg.MaxConnNum == 0 {
		cfg.MaxConnNum = 10000
	}
	if cfg.WriteBufferSize == 0 {
		cfg.WriteBufferSize = 4096
	}
	if cfg.HandshakeTimeout == 0 {
		cfg.HandshakeTimeout = 5
	}

	return &WsServer{
		addr:             cfg.Addr,
		wsMaxConnNum:     cfg.MaxConnNum,
		writeBufferSize:  cfg.WriteBufferSize,
		handshakeTimeout: time.Duration(cfg.HandshakeTimeout) * time.Second,
		clientPool: sync.Pool{
			New: func() any {
				return new(Client)
			},
		},
		registerChan:    make(chan *Client, 1000),
		unregisterChan:  make(chan *Client, 1000),
		kickHandlerChan: make(chan *kickHandler, 1000),
		validate:        validator.New(),
		Clients:         newUserMap(),
		// subscription: newSubscription(),
		Compressor:     NewGzipCompressor(),
		MessageHandler: NewServiceHandler(service.NewMessageService(database.GetDB()), producer),
	}
}

func (ws *WsServer) Run(ctx context.Context) {
	var client *Client

	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case client = <-ws.registerChan:
				ws.registerClient(client)
			case client = <-ws.unregisterChan:
				ws.unregisterClient(client)
			case onlineInfo := <-ws.kickHandlerChan:
				ws.multiTerminalLoginChecker(onlineInfo.clientOK, onlineInfo.oldClients, onlineInfo.newClient)
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		wsServer := http.Server{Addr: ws.addr, Handler: nil}
		http.HandleFunc("/ws", ws.wsHandler)
		go func() {
			defer close(done)
			<-ctx.Done()
			_ = wsServer.Shutdown(context.Background())
		}()
		log.Printf("WebSocket server starting on %s", ws.addr)
		err := wsServer.ListenAndServe()
		if err != nil {
			log.Printf("WebSocket server error: %v", err)
		}
		cancel(fmt.Errorf("msg gateway %w", err))
	}()

	<-ctx.Done()

}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {

	/*
		1. check max connection
		2. ckeck token
		3. upgrade to websocket
		4. create client
		5. register client
		6. start readMessage loop
	*/
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
		return
	}

	if _, err := util.ParseToken(r.URL.Query().Get(Token)); err != nil {
		log.Println("invalid token:", err)
		http.Error(w, "invalid token:"+err.Error(), http.StatusUnauthorized)
		return
	}

	upgrader := &websocket.Upgrader{
		HandshakeTimeout: ws.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
		WriteBufferSize:  ws.writeBufferSize,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade connection:"+err.Error(), http.StatusInternalServerError)
		return
	}

	client := ws.clientPool.Get().(*Client)
	client.ResetClient(w, r, conn, ws)
	ws.registerChan <- client

	go client.readMessage()
}

func (ws *WsServer) registerClient(client *Client) {
	oldClients, userOK, clientOK := ws.Clients.Get(client.UserID, client.PlatformID)

	if !userOK {
		ws.Clients.Set(client.UserID, client)
		ws.onlineUserNum.Add(1)
		ws.onlineUserConnNum.Add(1)
		prommetrics.OnlineUserGauge.Add(1)
	} else {
		ws.multiTerminalLoginChecker(clientOK, oldClients, client)

		ws.Clients.Set(client.UserID, client)
		ws.onlineUserConnNum.Add(1)
	}
}

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.Clients.DeleteClients(client.UserID, []*Client{client})
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
		prommetrics.OnlineUserGauge.Dec()
	}
	ws.onlineUserConnNum.Add(-1)
	// ws.subscription.DelClient(client)

}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

func (ws *WsServer) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	// kickTokenFunc := func(kickClients []*Client) {
	// 	var kickTokens []string
	// 	for _, c := range kickClients {
	// 		kickTokens = append(kickTokens, c.token)
	// 		c.KickOnlineMessage()
	// 	}
	// }
	// checkSameTokenFunc := func(oldClients []*Client) []*Client {
	// 	return nil
	// }
}
