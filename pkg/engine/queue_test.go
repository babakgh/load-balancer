package engine_test

import (
	"math/rand/v2"
	"reflect"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/babakgh/infillion/load-balancer/pkg/engine"
)

func TestQueue(t *testing.T) {
	totalItems := 1000
	queue := NewQueue[int]()

	got := make([]int, totalItems)
	want := make([]int, totalItems)

	for i := 0; i < totalItems; i++ {
		want[i] = rand.IntN(9999)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go produce(queue, &wg, want)
	go consume(queue, &wg, got)
	wg.Wait()

	queue.Close()

	assert.Equal(t, queue.Length(), 0)
	assert.Equal(t, got, want)
}

func produce(q *Queue[int], wg *sync.WaitGroup, items []int) {
	for i := 0; i < len(items); i++ {
		q.Enqueue(items[i])
	}
	wg.Done()
}

var result int

func consume(q *Queue[int], wg *sync.WaitGroup, items []int) {
	for i := 0; i < len(items); i++ {
		item, _ := q.Dequeue()
		items[i] = item
	}
	wg.Done()
}

func assertEqual[K comparable](t *testing.T, got, want K) {
	t.Helper()
	if got != want {
		t.Errorf("incorrect, got: %v, want %v", got, want)
	}
}

func assertDeepEqual[T any](t *testing.T, got, want []T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

// BenchmarkEnqueue tests the throughput of the Enqueue operation.
func BenchmarkEnqueue(b *testing.B) {
	queue := NewQueue[string]() // Assuming the queue is for strings

	b.ResetTimer() // Start the timer for the benchmark
	for i := 0; i < b.N; i++ {
		queue.Enqueue("item" + strconv.Itoa(i)) // Enqueue items as fast as possible
	}
}

// BenchmarkDequeue tests the throughput of the Dequeue operation.
func BenchmarkDequeue(b *testing.B) {
	queue := NewQueue[string]() // Assuming the queue is for strings

	// Preload the queue with more items than we plan to dequeue
	for i := 0; i < b.N+1000; i++ {
		queue.Enqueue("item" + strconv.Itoa(i))
	}

	b.ResetTimer() // Start the timer for the benchmark
	for i := 0; i < b.N; i++ {
		_, _ = queue.Dequeue() // Dequeue items as fast as possible
	}
}
