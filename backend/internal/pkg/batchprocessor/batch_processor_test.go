package batchprocessor

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestBasicEnqueueAndProcess verifies that items enqueued are processed by workers.
func TestBasicEnqueueAndProcess(t *testing.T) {
	bp := NewBatchProcessor[int]()

	var processed atomic.Int64
	var mu sync.Mutex
	results := make([]int, 0)

	// Update Do signature to accept []int
	bp.Do = func(ctx context.Context, channelID int, msgs []int) {
		processed.Add(int64(len(msgs)))
		mu.Lock()
		results = append(results, msgs...)
		mu.Unlock()
	}

	// Key function is required in current implementation
	bp.Key = func(data int) string {
		return fmt.Sprintf("%d", data%10)
	}

	// Start processor in background
	done := make(chan struct{})
	go func() {
		bp.Start()
		close(done)
	}()
	count := 1000
	// Enqueue items
	for i := 0; i < count; i++ {
		bp.Enqueue(i)
	}

	// Close to signal completion
	bp.Close()

	// Wait for processor to finish
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for processor to finish")
	}

	if processed.Load() != int64(count) {
		t.Errorf("expected %d processed, got %d", count, processed.Load())
	}
}

// TestTickerFlush verifies that data is flushed on ticker even if batch size not reached.
func TestTickerFlush(t *testing.T) {
	bp := NewBatchProcessor[int]()
	bp.size = 1000                       // large batch size
	bp.duration = 100 * time.Millisecond // short ticker

	var processed atomic.Int64

	bp.Do = func(ctx context.Context, channelID int, msgs []int) {
		processed.Add(int64(len(msgs)))
	}

	bp.Key = func(data int) string {
		return "same"
	}

	done := make(chan struct{})
	go func() {
		bp.Start()
		close(done)
	}()

	// Enqueue fewer items than batch size
	for i := 0; i < 10; i++ {
		bp.Enqueue(i)
	}

	// Wait for ticker to flush (should be < 200ms)
	time.Sleep(300 * time.Millisecond)

	// Close and wait
	bp.Close()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for processor to finish")
	}

	if processed.Load() != 10 {
		t.Errorf("expected 10 processed after ticker flush, got %d", processed.Load())
	}
}

// TestCloseWithoutEnqueue verifies that closing without enqueue doesn't hang.
func TestCloseWithoutEnqueue(t *testing.T) {
	bp := NewBatchProcessor[int]()

	bp.Do = func(ctx context.Context, channelID int, msgs []int) {
		// no-op
	}
	// Even if no data, Key should be defined to avoid potential nil pointer if logic changes
	bp.Key = func(i int) string { return "" }

	done := make(chan struct{})
	go func() {
		bp.Start()
		close(done)
	}()

	// Close immediately
	bp.Close()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout: processor should exit quickly when closed without data")
	}
}

// TestConsistentRouting verifies that same key always goes to same worker.
func TestConsistentRouting(t *testing.T) {
	bp := NewBatchProcessor[string]()

	workerForKey := make(map[string]int)
	var mu sync.Mutex

	bp.Do = func(ctx context.Context, channelID int, msgs []string) {
		mu.Lock()
		defer mu.Unlock()
		// All messages in a batch should have the same key if we group by key
		// But the current implementation groups by key in `schedule` map[string][]T
		// So `msgs` passed to Do will all share the same Key.

		// However, we need to verify that for a given key, it always goes to the same channelID
		if len(msgs) == 0 {
			return
		}

		// Check consistency within the batch (implicit by implementation, but good to check)
		firstKey := msgs[0] // Since Key func is identity

		if prev, ok := workerForKey[firstKey]; ok {
			if prev != channelID {
				t.Errorf("key %q routed to different workers: %d and %d", firstKey, prev, channelID)
			}
		} else {
			workerForKey[firstKey] = channelID
		}

		for _, val := range msgs {
			if val != firstKey {
				// This might happen if hash collision happens?
				// No, schedule groups by Key string. So all msgs in this slice MUST have same Key.
				t.Errorf("Batch contained mixed keys: %s and %s", firstKey, val)
			}
		}
	}

	bp.Key = func(data string) string {
		return data
	}

	done := make(chan struct{})
	go func() {
		bp.Start()
		close(done)
	}()

	// Enqueue same keys multiple times
	keys := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for i := 0; i < 50; i++ {
		for _, k := range keys {
			bp.Enqueue(k)
		}
	}

	bp.Close()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for processor to finish")
	}
}

// TestBatchingLogic verifies that items are actually batched together
func TestBatchingLogic(t *testing.T) {
	bp := NewBatchProcessor[int]()
	bp.size = 10
	bp.duration = 1 * time.Second // Long duration to force size-based batching

	var batchCount atomic.Int64
	var itemCount atomic.Int64

	bp.Do = func(ctx context.Context, channelID int, msgs []int) {
		batchCount.Add(1)
		itemCount.Add(int64(len(msgs)))
	}

	bp.Key = func(i int) string { return "same-key" }

	done := make(chan struct{})
	go func() {
		bp.Start()
		close(done)
	}()

	// Send 10 items (equal to batch size)
	for i := 0; i < 10; i++ {
		bp.Enqueue(i)
	}

	// Wait a bit for processing (but less than duration)
	time.Sleep(100 * time.Millisecond)

	bp.Close()
	<-done

	if itemCount.Load() != 10 {
		t.Errorf("Expected 10 items, got %d", itemCount.Load())
	}

	// Should ideally be 1 batch if processed fast enough
	if batchCount.Load() != 1 {
		t.Errorf("Expected 1 batch, got %d", batchCount.Load())
	}
}
