package engine

import (
	"fmt"
	"sync"
	"time"
)

type ReportManager struct {
	meters   sync.Map
	ticker   *time.Ticker
	interval int64
}

func NewReportManager(interval int64) *ReportManager {
	rm := &ReportManager{
		ticker:   time.NewTicker(time.Duration(interval) * time.Second),
		interval: interval,
	}
	go rm.run()
	return rm
}

func (rm *ReportManager) Register(id string, windowLengthMinutes int64) {
	alpha := rm.calculateAlpha(windowLengthMinutes)
	rm.meters.Store(id, NewRateMeter(alpha, rm.interval))
}

func (rm *ReportManager) run() {
	for range rm.ticker.C {
		rm.meters.Range(func(key, value interface{}) bool {
			value.(*RateMeter).Tick()
			// meter := value.(*RateMeter)
			// meter.Rate()
			return true
		})
	}
}

func (rm *ReportManager) Increment(id string) {
	if v, ok := rm.meters.Load(id); ok {
		v.(*RateMeter).Update(1)
	}
}

func (rm *ReportManager) Report() {
	rm.meters.Range(func(key, value interface{}) bool {
		meter := value.(*RateMeter)
		rate := meter.Rate()
		fmt.Printf("%v: %.3f -> %v\n", key, rate, rate)
		return true
	})
}

func (rm *ReportManager) Stop() {
	rm.ticker.Stop()
}

func (rm *ReportManager) calculateAlpha(window int64) float64 {
	// more reponsive to changes in interval
	// alpha := 1 - math.Exp(float64(-interval)/float64(window*60))

	N := window * 60 / rm.interval
	alpha := 2.0 / float64(N+1)
	return alpha
}
