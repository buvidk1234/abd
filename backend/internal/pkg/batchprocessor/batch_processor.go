package batchprocessor

import (
	"context"
	"hash/fnv"
	"log"
	"sync"
	"time"
)

const (
	DefaultSize        = 1024
	DefaultChanSize    = 1024
	DefaultDuration    = time.Second
	DefaultWorkerCount = 5
)

type BatchProcessor[T any] struct {
	size        int
	duration    time.Duration
	workerCount int
	data        chan T
	subData     []chan []T
	Key         func(data T) string

	OnComplete func(lastMessage *T, totalCount int)

	Do     func(ctx context.Context, channelID int, msgs []T)
	waiter sync.WaitGroup
}

func NewBatchProcessor[T any]() *BatchProcessor[T] {

	subData := make([]chan []T, DefaultWorkerCount)
	for i := 0; i < DefaultWorkerCount; i++ {
		subData[i] = make(chan []T, DefaultChanSize)
	}

	return &BatchProcessor[T]{
		size:        DefaultSize,
		duration:    DefaultDuration,
		workerCount: DefaultWorkerCount,
		data:        make(chan T, DefaultChanSize),
		// preallocate subdata slice with capacity equal to worker count
		subData: subData,
	}
}

func (bp *BatchProcessor[T]) Start() {
	// start workers
	bp.waiter.Add(bp.workerCount)
	for i := 0; i < bp.workerCount; i++ {
		go bp.work(i)
	}

	// start scheduler
	bp.waiter.Add(1)
	go bp.schedule()

	// wait for scheduler and workers to finish
	bp.waiter.Wait()
}

func (bp *BatchProcessor[T]) schedule() {
	ticker := time.NewTicker(bp.duration)

	defer func() {
		ticker.Stop()
		// close worker channels when scheduler exits
		for i := 0; i < bp.workerCount; i++ {
			close(bp.subData[i])
		}

		bp.waiter.Done()
	}()

	counter := 0
	dataTmp := make(map[string][]T)
	for {
		select {
		case v, ok := <-bp.data:
			if !ok {
				log.Println("batch processor has been closed")
				// input closed: dispatch any remaining data, then return
				if counter > 0 {
					bp.dispatch(dataTmp, counter)
				}
				return
			}

			dataTmp[bp.Key(v)] = append(dataTmp[bp.Key(v)], v)
			counter++
			if counter >= bp.size {
				log.Println("batch size reached")
				// dispatch data
				bp.dispatch(dataTmp, counter)
				// reset counter and reuse slice
				counter = 0
				dataTmp = make(map[string][]T)
			}
		case <-ticker.C:
			// log.Println("ticker ticked")
			if counter > 0 {
				// dispatch data
				bp.dispatch(dataTmp, counter)
				// reset counter and reuse slice
				counter = 0
				dataTmp = make(map[string][]T)
			}
		}
	}
}

func (bp *BatchProcessor[T]) dispatch(data map[string][]T, count int) {
	rr := 0
	for _, v := range data {
		var idx int
		if bp.Key == nil {
			// round-robin when no Key function provided
			idx = rr % bp.workerCount
			rr++
		} else {
			idx = bp.getIdx(v[0])
		}
		bp.subData[idx] <- v
	}
}

func (bp *BatchProcessor[T]) getIdx(data T) int {

	key := bp.Key(data)
	// hash key to pick worker
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	idx := int(h.Sum32() % uint32(bp.workerCount))
	return idx
}

func (bp *BatchProcessor[T]) work(i int) {
	defer bp.waiter.Done()
	// process items until subdata channel is closed
	for v := range bp.subData[i] {
		if bp.Do != nil {
			// execute Do; pass background context for now
			bp.Do(context.Background(), i, v)
		}
	}
}

// Enqueue adds an item to the processor. Returns false if the input channel is closed.
func (bp *BatchProcessor[T]) Enqueue(item T) bool {
	defer func() { recover() }() // recover from send on closed channel
	bp.data <- item
	return true
}

// Close closes the input channel, signaling the scheduler to finish processing and exit.
func (bp *BatchProcessor[T]) Close() {
	close(bp.data)
}
