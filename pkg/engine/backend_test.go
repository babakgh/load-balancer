package engine_test

import (
	"time"

	. "github.com/babakgh/load-balancer/pkg/engine"
)

// MockBackend is a simple implementation of Backender
type MockBackend struct {
	IDVal  string
	KeyVal string
}

func (m *MockBackend) ID() string {
	return m.IDVal
}

func (m *MockBackend) Key() string {
	return m.KeyVal
}

func (m *MockBackend) Process(request Requester) error {
	// Do some work
	time.Sleep(30 * time.Millisecond)
	return nil
}
