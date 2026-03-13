package main

import (
	"backend/tests/stress/common"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	apiAddr     = flag.String("api", "http://127.0.0.1:8080", "API server address")
	wsAddr      = flag.String("ws", "ws://127.0.0.1:8082/ws", "WebSocket server address")
	numMsg      = flag.Int("n", 100, "Number of messages to send")
	interval    = flag.Int("i", 500, "Interval between messages (ms)")
	senderIdx   = flag.Int("sender", 10001, "Sender username index")
	receiverIdx = flag.Int("receiver", 10002, "Receiver username index")
	logFile     = flag.String("log", "reliability_test.log", "Log file path")
)

type MsgTracker struct {
	SendTime map[string]time.Time
	RecvTime map[string]time.Time
	mu       sync.Mutex
}

func main() {
	flag.Parse()

	log.Printf("Starting reliability test: sender index %d, receiver index %d, %d messages",
		*senderIdx, *receiverIdx, *numMsg)

	// 1. Auth users
	senderID, senderToken, err := authenticate(fmt.Sprintf("rel_%d", *senderIdx))
	if err != nil {
		log.Fatal(err)
	}
	receiverID, receiverToken, err := authenticate(fmt.Sprintf("rel_%d", *receiverIdx))
	if err != nil {
		log.Fatal(err)
	}

	tracker := &MsgTracker{
		SendTime: make(map[string]time.Time),
		RecvTime: make(map[string]time.Time),
	}

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	var wg sync.WaitGroup
	wg.Add(2)

	// 2. Receiver start
	go func() {
		defer wg.Done()
		receiveLoop(receiverToken, tracker)
	}()

	// Wait a bit for receiver to be ready
	time.Sleep(2 * time.Second)

	// 3. Sender start
	go func() {
		defer wg.Done()
		sendLoop(senderID, senderToken, receiverID, tracker)
	}()

	wg.Wait()

	// 4. Report
	analyze(tracker)
}

func authenticate(username string) (int64, string, error) {
	password := "123456"
	token, err := common.Login(*apiAddr, username, password)
	if err != nil {
		_, err := common.Register(*apiAddr, username, password)
		if err != nil {
			return 0, "", err
		}
		token, err = common.Login(*apiAddr, username, password)
		if err != nil {
			return 0, "", err
		}
	}
	id, err := common.GetInfo(*apiAddr, token)
	return id, token, err
}

func dialWithRetry(url string, maxRetries int) (*websocket.Conn, error) {
	var conn *websocket.Conn
	var err error
	for i := 0; i < maxRetries; i++ {
		dialer := websocket.DefaultDialer
		dialer.HandshakeTimeout = 10 * time.Second
		conn, _, err = dialer.Dial(url, nil)
		if err == nil {
			return conn, nil
		}
		log.Printf("Dial attempt %d failed: %v, retrying in 2s...", i+1, err)
		time.Sleep(2 * time.Second)
	}
	return nil, err
}

func sendLoop(senderID int64, token string, receiverID int64, tracker *MsgTracker) {
	wsURL := common.GetWSURL(*wsAddr, token, common.WebPlatformID)
	conn, err := dialWithRetry(wsURL, 5)
	if err != nil {
		log.Fatalf("Sender failed to connect after retries: %v", err)
	}
	defer conn.Close()
	log.Printf("Sender connected successfully")

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
	for i := 0; i < *numMsg; i++ {
		clientMsgID := fmt.Sprintf("rel_msg_%d_%d", senderID, time.Now().UnixNano())

		msgReq := map[string]any{
			"sender_id":     strconv.FormatInt(senderID, 10),
			"conv_type":     1, // Single chat
			"target_id":     strconv.FormatInt(receiverID, 10),
			"msg_type":      1, // Text
			"client_msg_id": clientMsgID,
			"content":       "reliability test " + clientMsgID,
		}
		msgData, _ := json.Marshal(msgReq)

		req := common.InboundReq{
			ReqIdentifier: common.WSSendMsg,
			MsgIncr:       fmt.Sprintf("%d", i),
			Data:          msgData,
		}

		payload, _ := json.Marshal(req)

		tracker.mu.Lock()
		tracker.SendTime[clientMsgID] = time.Now()
		tracker.mu.Unlock()

		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("Send error: %v", err)
		}

		time.Sleep(time.Duration(*interval) * time.Millisecond)
	}
	// Give some time for the last message to arrive
	time.Sleep(2 * time.Second)
}

func receiveLoop(token string, tracker *MsgTracker) {
	wsURL := common.GetWSURL(*wsAddr, token, common.WebPlatformID)
	conn, err := dialWithRetry(wsURL, 5)
	if err != nil {
		log.Fatalf("Receiver failed to connect after retries: %v", err)
	}
	defer conn.Close()
	log.Printf("Receiver connected successfully")

	for {
		_ = conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Receiver read error: %v", err)
			return
		}
		var resp common.Resp
		if err := json.Unmarshal(msg, &resp); err != nil {
			log.Printf("Receiver unknown msg: %s", string(msg))
			continue
		}

		if resp.ReqIdentifier != common.WSPushMsg {
			// log.Printf("Receiver got other msg: ID=%d", resp.ReqIdentifier)
			continue
		}

		// log.Printf("Receiver got Push: %s", string(resp.Data))
		// Extract client_msg_id from pushed data
		var pushedMsg struct {
			ClientMsgID string `json:"client_msg_id"`
		}
		json.Unmarshal(resp.Data, &pushedMsg)

		if pushedMsg.ClientMsgID != "" {
			tracker.mu.Lock()
			if _, ok := tracker.SendTime[pushedMsg.ClientMsgID]; ok {
				if _, already := tracker.RecvTime[pushedMsg.ClientMsgID]; !already {
					tracker.RecvTime[pushedMsg.ClientMsgID] = time.Now()
					// log.Printf("Message %s matched and tracked", pushedMsg.ClientMsgID)
				}
			} else {
				// log.Printf("Received message %s but not in sent map", pushedMsg.ClientMsgID)
			}
			tracker.mu.Unlock()
		}

		tracker.mu.Lock()
		recvCount := len(tracker.RecvTime)
		tracker.mu.Unlock()
		if recvCount >= *numMsg {
			log.Printf("Received all %d messages, exiting", *numMsg)
			return
		}
	}
}

func analyze(tracker *MsgTracker) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	var latencies []time.Duration
	received := 0

	for id, sendTime := range tracker.SendTime {
		if recvTime, ok := tracker.RecvTime[id]; ok {
			latencies = append(latencies, recvTime.Sub(sendTime))
			received++
		}
	}

	log.Printf("\n--- Reliability Test Result ---")
	log.Printf("Messages Sent:     %d", len(tracker.SendTime))
	log.Printf("Messages Received: %d", received)
	log.Printf("Reliability:       %.2f%%", float64(received)/float64(len(tracker.SendTime))*100)

	if received > 0 {
		sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

		var total time.Duration
		for _, l := range latencies {
			total += l
		}

		avg := total / time.Duration(received)
		p95 := latencies[int(float64(received)*0.95)]
		p99Idx := int(float64(received) * 0.99)
		if p99Idx >= received {
			p99Idx = received - 1
		}
		p99 := latencies[p99Idx]

		log.Printf("Avg Latency:       %v", avg)
		log.Printf("Max Latency:       %v", latencies[received-1])
		log.Printf("Min Latency:       %v", latencies[0])
		log.Printf("P95 Latency:       %v", p95)
		log.Printf("P99 Latency:       %v", p99)
	}
	log.Printf("--------------------------------")
}
