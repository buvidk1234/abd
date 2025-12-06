package im

import (
	"backend/internal/pkg/database"
	"backend/internal/pkg/kafka"
	"backend/internal/service"
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type WsServer struct {
	port              int
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
}

type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func NewWsServer() *WsServer {
	producer, err := kafka.NewSyncProducer()
	if err != nil {
		panic(fmt.Sprintf("failed to create kafka producer: %v", err))
	}
	return &WsServer{
		port:             8081,
		wsMaxConnNum:     10000,
		writeBufferSize:  4096,
		handshakeTimeout: 5 * time.Second,
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
		wsServer := http.Server{Addr: fmt.Sprintf(":%d", ws.port), Handler: nil}
		http.HandleFunc("/ws", ws.wsHandler)
		go func() {
			defer close(done)
			<-ctx.Done()
			_ = wsServer.Shutdown(context.Background())
		}()
		err := wsServer.ListenAndServe()
		if err == nil {
			err = fmt.Errorf("http server closed")
		}
		cancel(fmt.Errorf("msg gateway %w", err))
	}()

	<-ctx.Done()

}

func (ws *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {

	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		http.Error(w, "too many connections", http.StatusServiceUnavailable)
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
