package engine

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// Exponential moving average meter (EMA)
type RateMeter struct {
	uncounted int64      // Number of events that have occurred
	alpha     float64    // Smoothing factor for the exponential moving average
	rate      uint64     // Current rate, stored as bits for atomic operations
	init      uint32     // Initialization check flag
	mutex     sync.Mutex // Mutex for initializing the ticker
	interval  int64      // Interval in seconds for EMA rate calculation
	// window    int64      // time window in minutes for EMA rate calculation
}

// NewRateMeter initializes a RateMeter with a specified alpha and interval in seconds and starts the ticker.
func NewRateMeter(alpha float64, interval int64) *RateMeter {
	rm := &RateMeter{
		alpha:    alpha,
		interval: interval,
		// window:   window,
	}
	return rm
}

func (c *RateMeter) Tick() {
	if atomic.LoadUint32(&c.init) == 1 { // hot path
		c.updateRate(c.fetchInstantRate())
	} else { // cold path
		c.mutex.Lock()
		if atomic.LoadUint32(&c.init) == 1 {
			c.updateRate(c.fetchInstantRate())
		} else {
			atomic.StoreUint32(&c.init, 1)
			atomic.StoreUint64(&c.rate, math.Float64bits(c.fetchInstantRate()))
		}
		c.mutex.Unlock()
	}
}

func (m *RateMeter) Update(n int64) {
	atomic.AddInt64(&m.uncounted, n)
}

func (m *RateMeter) Rate() float64 {
	currentRateBits := atomic.LoadUint64(&m.rate)
	return math.Float64frombits(currentRateBits)
}

func (m *RateMeter) fetchInstantRate() float64 {
	count := atomic.SwapInt64(&m.uncounted, 0)
	instantRate := float64(count) / (float64(m.interval) * float64(time.Second))
	return instantRate
}

func (m *RateMeter) updateRate(instantRate float64) {
	currentRate := math.Float64frombits(atomic.LoadUint64(&m.rate))
	currentRate += m.alpha * (instantRate - currentRate)
	atomic.StoreUint64(&m.rate, math.Float64bits(currentRate))
}
