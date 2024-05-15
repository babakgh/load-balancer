package engine

import "sync/atomic"

type RoundRobinSelector struct {
	Total int64
	index int64
}

func NewRoundRobinSelector(total int64) *RoundRobinSelector {
	return &RoundRobinSelector{
		Total: total,
		index: 0,
	}
}

func (r *RoundRobinSelector) Select(_ Requester) int64 {
	current := atomic.LoadInt64(&r.index)
	next := (current + 1) % r.Total
	atomic.StoreInt64(&r.index, next)
	return current
}
