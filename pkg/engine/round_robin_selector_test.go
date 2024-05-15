package engine_test

import (
	"testing"

	. "github.com/babakgh/infillion/load-balancer/pkg/engine"
	"github.com/stretchr/testify/assert"
)

func TestRoundRobinSelector(t *testing.T) {
	// Create mock backends
	selector := RoundRobinSelector{Total: 3}

	// Simulate requests and check for correct backend selection
	want := []int64{0, 1, 2, 0, 1, 2, 0, 1, 2, 1}
	got := []int64{}
	for i := 0; i < 10; i++ {
		got = append(got, selector.Select(nil))
	}

	assert.Equal(t, want, got, "Test failed at iteration %d: expected backend %v, got %v", want, got)

}
