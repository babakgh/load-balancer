package engine

import (
	"sync"
)

// Queue represents a thread-safe FIFO (first in, first out) structure that can hold items of any type.
type Queue[T any] struct {
	items   []T           // items hold the elements of the queue.
	lock    sync.Mutex    // lock protects access to the queue.
	state   *sync.Cond    // state is used to signal changes in the queue's state.
	closeCh chan struct{} // closeCh is used to signal that the queue is closing.
	closed  bool          // closed indicates whether the queue has been closed.
}

// NewQueue creates and returns a new Queue instance.
func NewQueue[T any]() *Queue[T] {
	q := &Queue[T]{
		closeCh: make(chan struct{}),
	}
	q.state = sync.NewCond(&q.lock) // Initialize the condition variable.
	return q
}

// Enqueue adds an item to the end of the queue and signals any waiting goroutines.
func (q *Queue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.items = append(q.items, item)
	q.state.Signal() // Notify goroutines waiting in Dequeue that the state has changed.
}

// Dequeue removes and returns the front item from the queue.
// If the queue is empty, it blocks until an item is available or the queue is closed.
// Returns false if the queue is closed and no items are available.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for len(q.items) == 0 {
		select {
		case <-q.closeCh:
			return *new(T), false // Return zero value and false if the queue is closed.
		default:
			q.state.Wait() // Wait for an item to be enqueued or the queue to close.
		}
	}

	item := q.items[0]
	q.items = q.items[1:]

	return item, true
}

// Close marks the queue as closed and wakes up all waiting goroutines.
func (q *Queue[T]) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if !q.closed {
		q.closed = true
		close(q.closeCh)    // Close the channel to signal closure.
		q.state.Broadcast() // Wake up all goroutines waiting on the condition variable.
	}
}

// Length returns the current number of items in the queue.
func (q *Queue[T]) Length() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.items)
}
