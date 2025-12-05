package im

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWsHandler_RegisterClient(t *testing.T) {
	ws := NewWsServer()

	srv := httptest.NewServer(http.HandlerFunc(ws.wsHandler))
	defer srv.Close()

	// build ws url from httptest server url
	u := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws?platform_id=3&token=abc"

	d := websocket.Dialer{}
	conn, _, err := d.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.Close()

	select {
	case client := <-ws.registerChan:
		if client == nil {
			t.Fatalf("received nil client")
		}
		if client.PlatformID != 3 {
			t.Fatalf("expected platform 3, got %d", client.PlatformID)
		}
		if client.token != "abc" {
			t.Fatalf("expected token 'abc', got '%s'", client.token)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timed out waiting for client registration")
	}
}

func TestWsServer_Run_RegisterAndShutdown(t *testing.T) {
	ws := NewWsServer()

	// pick a free port
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	addr := l.Addr().(*net.TCPAddr)
	port := addr.Port
	_ = l.Close()

	ws.port = port

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run server
	go ws.Run(ctx)

	// allow server to start
	time.Sleep(100 * time.Millisecond)

	u := "ws://127.0.0.1:" + fmt.Sprintf("%d", port) + "/ws?platform_id=7&token=tok-run"
	d := websocket.Dialer{}
	conn, _, err := d.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	defer conn.Close()

	// cancel context to stop server
	cancel()

	// give server a moment to shutdown
	time.Sleep(200 * time.Millisecond)
}
