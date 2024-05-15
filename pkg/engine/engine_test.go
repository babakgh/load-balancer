package engine_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	. "github.com/babakgh/infillion/load-balancer/pkg/engine"
	"github.com/stretchr/testify/assert"
)

func TestEngine(t *testing.T) {
	runningTimeInMinute := 10
	requestsPerSecond := 60000
	maxAvailableBackends := 100
	maxEngineWorkers := 2500

	averageRequestPerSecondPerBackend := requestsPerSecond / maxAvailableBackends
	t.Logf("averageRequestPerSecondPerBackend: %v\n", averageRequestPerSecondPerBackend)

	engine := createEngine(maxEngineWorkers, maxAvailableBackends)
	t.Logf("maxEngineWorkers: %v\n", maxEngineWorkers)

	queue := createQueue()
	t.Logf("queue is created\n")

	backends := createAllBackends(maxAvailableBackends)
	t.Logf("maxAvailableBackends: %v\n", maxAvailableBackends)
	// last_key := backends[len(backends)-1].Key()

	t.Logf("Enqueueing %d requests per second for %d minutes\n", requestsPerSecond, runningTimeInMinute)
	var wg sync.WaitGroup
	wg.Add(1)

	go enqueueRequestPerSecond(queue, runningTimeInMinute, requestsPerSecond)
	go SetTimeout(runningTimeInMinute, &wg)

	t.Logf("Engine is starting ... \n")
	engine.Start(queue, backends)

	t.Logf("Engine is started at %v\n", time.Now())
	wg.Wait()

	t.Logf("Engine is stopping ...\n")
	engine.Stop()
	t.Logf("Engine is stopped! at %v\n", time.Now())

	got := queue.Length()
	want := 0
	assert.Equal(t, want, got, "Not all items were processed, %d items left", got)
}

// func BenchmarkEngine(b *testing.B) {
// }

func createAllBackends(maxAvailableBackends int) []Backender {
	backends := make([]Backender, maxAvailableBackends)
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	total := len(letters)

	for i := 0; i <= maxAvailableBackends-1; i++ {
		key := string(letters[i/(total*total)%total]) +
			string(letters[i/total%total]) +
			string(letters[i%total])
		backends[i] = &MockBackend{IDVal: strconv.Itoa(i), KeyVal: key}
	}
	return backends
}

func createQueue() *Queue[Requester] {
	return NewQueue[Requester]()
}

func createEngine(maxEngineWorkers int, maxAvailableBackends int) *Engine {
	backendSelector := RoundRobinSelector{Total: int64(maxAvailableBackends)}
	engine := NewEngine(uint32(maxEngineWorkers), &backendSelector)
	return engine
}

func enqueueRequestPerSecond(queue *Queue[Requester], runningTimeInMinute int, requestsPerSecond int) {
	for n := 0; n < runningTimeInMinute*60; n++ {
		now := time.Now()
		for i := 1; i <= requestsPerSecond; i++ {
			id := strconv.FormatUint(uint64(i), 10)
			queue.Enqueue(&MockRequest{IDVal: id, KeyVal: "key" + id})
		}
		// change this to t.Logf
		fmt.Printf("Enqueued %d items\n", requestsPerSecond)
		time.Sleep(1000*time.Millisecond - time.Duration(time.Since(now).Milliseconds()))
	}
}

func SetTimeout(runningTimeInMinute int, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(runningTimeInMinute) * 60 * time.Second)
	wg.Done()
}
