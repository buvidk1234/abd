package im

import (
	"backend/internal/pkg/kafka"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClient_ResetClient(t *testing.T) {
	// 初始化 Kafka 配置，防止 NewWsServer panic
	kafka.Init(kafka.Config{Addr: []string{"192.168.6.130:9092"}})

	// Setup a test server to handle the websocket connection
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
	}))
	defer s.Close()

	// Connect to the test server
	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/?sendID=user1&platformID=1&compression=gzip&token=test_token"
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// Create a dummy request and response writer for ResetClient
	req := httptest.NewRequest("GET", u, nil)
	w := httptest.NewRecorder()

	client := &Client{}
	wsServer := NewWsServer() // Assuming NewWsServer is available and simple enough

	client.ResetClient(w, req, conn, wsServer)

	if client.UserID != "user1" {
		t.Errorf("expected UserID 'user1', got '%s'", client.UserID)
	}
	if client.PlatformID != 1 {
		t.Errorf("expected PlatformID 1, got %d", client.PlatformID)
	}
	if !client.IsCompress {
		t.Errorf("expected IsCompress true, got %v", client.IsCompress)
	}
	if client.token != "test_token" {
		t.Errorf("expected token 'test_token', got '%s'", client.token)
	}
	if client.server != wsServer {
		t.Errorf("expected server to be set")
	}
	if client.Encoder == nil {
		t.Errorf("expected Encoder to be set")
	}
}

func TestClient_PushUserOnlineStatus(t *testing.T) {
	// Setup a test server that reads messages
	var receivedMsg []byte
	done := make(chan struct{})

	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, msg, err := c.ReadMessage()
		if err != nil {
			return
		}
		receivedMsg = msg
		close(done)
	}))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	client := &Client{
		conn:    conn,
		Encoder: NewJsonEncoder(),
	}

	data := []byte("online")
	err = client.PushUserOnlineStatus(data)
	if err != nil {
		t.Fatalf("PushUserOnlineStatus failed: %v", err)
	}

	select {
	case <-done:
		// Verify the message content if possible.
		// The message sent is binary encoded Resp.
		// We can try to decode it if we know the format, or just assume success if we received something.
		if len(receivedMsg) == 0 {
			t.Error("received empty message")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}

func TestClient_KickOnlineMessage(t *testing.T) {
	// Setup a test server
	done := make(chan struct{})
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, _, err = c.ReadMessage()
		if err == nil {
			close(done)
		}
	}))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	// defer conn.Close() // KickOnlineMessage closes the connection

	// We need to set hbCtx for close() to work without panic if it cancels context
	ctx, cancel := context.WithCancel(context.Background())

	wsServer := NewWsServer()
	// Consume unregister channel to prevent blocking
	go func() {
		for range wsServer.unregisterChan {
		}
	}()

	client := &Client{
		conn:     conn,
		Encoder:  NewJsonEncoder(),
		hbCtx:    ctx,
		hbCancel: cancel,
		server:   wsServer,
	}

	err = client.KickOnlineMessage()
	if err != nil {
		t.Fatalf("KickOnlineMessage failed: %v", err)
	}

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for message")
	}

	// Verify connection is closed?
	// client.close() calls conn.Close()
}
