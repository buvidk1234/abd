package main

import (
	"backend/tests/stress/common"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	apiAddr       = flag.String("api", "http://127.0.0.1:8080", "API server address")
	wsAddr        = flag.String("ws", "ws://127.0.0.1:8082/ws", "WebSocket server address")
	numUsers      = flag.Int("u", 100, "Number of users to simulate")
	loginInterval = flag.Int("li", 100, "Interval between logins (ms)")
	msgInterval   = flag.Int("mi", 1000, "Average interval between messages per user (ms)")
	duration      = flag.Int("d", 60, "Test duration (seconds)")
	startIdx      = flag.Int("s", 1, "Starting index for usernames")
	logFile       = flag.String("log", "pressure_test.log", "Log file path")
)

type Stats struct {
	ActiveConns atomic.Int64
	SentMsgs    atomic.Int64
	RecvMsgs    atomic.Int64
	Errors      atomic.Int64
	Latencies   []time.Duration
	latMu       sync.Mutex
}

var stats Stats

type MessageInfo struct {
	ClientMsgID string
	SendTime    time.Time
}

var (
	msgInFlight sync.Map // clientMsgID -> time.Time
)

type UserInfo struct {
	ID    int64
	Token string
}

var (
	userList []UserInfo
	userMu   sync.Mutex
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, f))

	log.Printf("Starting pressure test: %d users, login interval %dms, msg interval %dms, duration %ds",
		*numUsers, *loginInterval, *msgInterval, *duration)

	wg := sync.WaitGroup{}
	stopChan := make(chan struct{})

	// Concurrent login/register all users
	log.Printf("Authenticating %d users with concurrency...", *numUsers)
	var authWg sync.WaitGroup
	authChan := make(chan int, 50) // Concurrency limit for auth
	
	for i := 0; i < *numUsers; i++ {
		authWg.Add(1)
		go func(i int) {
			defer authWg.Done()
			authChan <- 1
			defer func() { <-authChan }()

			idx := *startIdx + i
			username := fmt.Sprintf("stress_%d", idx)
			password := "123456"

			id, token, err := authenticate(username, password)
			if err != nil {
				log.Printf("User %s failed to auth: %v", username, err)
				stats.Errors.Add(1)
				return
			}
			
			userMu.Lock()
			userList = append(userList, UserInfo{ID: id, Token: token})
			count := len(userList)
			userMu.Unlock()
			
			if count%100 == 0 {
				log.Printf("Progress: Authenticated %d/%d users...", count, *numUsers)
			}
		}(i)
	}
	authWg.Wait()

	if len(userList) < 2 {
		log.Fatal("Not enough users authenticated for the test")
	}
	log.Printf("Starting simulation with %d authenticated users", len(userList))

	for _, user := range userList {
		wg.Add(1)
		go func(u UserInfo) {
			defer wg.Done()
			simulateUser(u, stopChan)
		}(user)

		if *loginInterval > 0 {
			time.Sleep(time.Duration(*loginInterval) * time.Millisecond)
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	startTime := time.Now()
	go func() {
		for range ticker.C {
			elapsed := time.Since(startTime).Seconds()
			log.Printf("Elapsed: %.0fs | Active: %d | Sent: %d | Recv: %d | Err: %d | TPS: %.1f",
				elapsed, stats.ActiveConns.Load(), stats.SentMsgs.Load(), stats.RecvMsgs.Load(), stats.Errors.Load(),
				float64(stats.SentMsgs.Load())/elapsed)
		}
	}()

	time.Sleep(time.Duration(*duration) * time.Second)
	close(stopChan)
	ticker.Stop()
	wg.Wait()

	finalReport(startTime)
}

func finalReport(startTime time.Time) {
	elapsed := time.Since(startTime).Seconds()
	stats.latMu.Lock()
	lats := make([]time.Duration, len(stats.Latencies))
	copy(lats, stats.Latencies)
	stats.latMu.Unlock()

	sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })

	var avg time.Duration
	if len(lats) > 0 {
		var total time.Duration
		for _, l := range lats {
			total += l
		}
		avg = total / time.Duration(len(lats))
	}

	p95 := time.Duration(0)
	p99 := time.Duration(0)
	if len(lats) > 0 {
		p95 = lats[int(float64(len(lats))*0.95)]
		p99 = lats[int(float64(len(lats))*0.99)]
	}

	log.Printf("\n--- Final Pressure Test Report ---")
	log.Printf("Test Duration:    %.2fs", elapsed)
	log.Printf("Total Sent:       %d", stats.SentMsgs.Load())
	log.Printf("Total Recv:       %d", stats.RecvMsgs.Load())
	log.Printf("Total Errors:     %d", stats.Errors.Load())
	log.Printf("Average TPS:      %.2f", float64(stats.SentMsgs.Load())/elapsed)
	if len(lats) > 0 {
		log.Printf("Avg Latency:      %v", avg)
		log.Printf("Max Latency:      %v", lats[len(lats)-1])
		log.Printf("Min Latency:      %v", lats[0])
		log.Printf("P95 Latency:      %v", p95)
		log.Printf("P99 Latency:      %v", p99)
	}
	log.Printf("----------------------------------\n")
}

func authenticate(username, password string) (int64, string, error) {
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

func simulateUser(user UserInfo, stop chan struct{}) {
	wsURL := common.GetWSURL(*wsAddr, user.Token, common.WebPlatformID)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		stats.Errors.Add(1)
		return
	}
	defer conn.Close()

	stats.ActiveConns.Add(1)
	defer stats.ActiveConns.Add(-1)

	// Read loop
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			stats.RecvMsgs.Add(1)

			var resp common.Resp
			if err := json.Unmarshal(msg, &resp); err == nil {
				// We expect push or ack. In ABD, sent message success often comes back as an ACK or push.
				// If we have a push with client_msg_id, we can calculate latency.
				if resp.ReqIdentifier == common.WSPushMsg {
					var pushedMsg struct {
						ClientMsgID string `json:"client_msg_id"`
					}
					json.Unmarshal(resp.Data, &pushedMsg)
					if pushedMsg.ClientMsgID != "" {
						if startTime, ok := msgInFlight.LoadAndDelete(pushedMsg.ClientMsgID); ok {
							latency := time.Since(startTime.(time.Time))
							stats.latMu.Lock()
							stats.Latencies = append(stats.Latencies, latency)
							stats.latMu.Unlock()
						}
					}
				}
			}
		}
	}()

	// Write loop
	for {
		select {
		case <-stop:
			return
		case <-time.After(time.Duration(*msgInterval/2+rand.Intn(*msgInterval)) * time.Millisecond):
			targetUser := userList[rand.Intn(len(userList))]
			// In ABD, ID 0 might be problematic if we use it for SendMessage.
			// Let's assume we need real IDs of users. 
			// I will update authenticate to return ID if possible.

			clientMsgID := fmt.Sprintf("%d_%d_%d", user.ID, targetUser.ID, time.Now().UnixNano())
			msgReq := map[string]any{
				"sender_id":     strconv.FormatInt(user.ID, 10),
				"conv_type":     1, // Single chat
				"target_id":     strconv.FormatInt(targetUser.ID, 10),
				"msg_type":      1, // Text
				"client_msg_id": clientMsgID,
				"content":       "pressure test message " + time.Now().String(),
			}
			msgData, _ := json.Marshal(msgReq)

			req := common.InboundReq{
				ReqIdentifier: common.WSSendMsg,
				MsgIncr:       fmt.Sprintf("%d", time.Now().UnixNano()),
				Data:          msgData,
			}

			payload, _ := json.Marshal(req)
			msgInFlight.Store(clientMsgID, time.Now())
			if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				stats.Errors.Add(1)
				return
			}
			stats.SentMsgs.Add(1)
		}
	}
}
