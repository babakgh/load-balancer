package engine

import (
	"log"
	"sync"
	"time"
)

type Engine struct {
	queue      *Queue[Requester]
	done       chan struct{}
	wg         sync.WaitGroup
	maxWorkers uint32

	reporter *ReportManager
	backends []Backender
	selector Selector
}

func NewEngine(maxWorkers uint32, selector Selector) *Engine {
	return &Engine{
		maxWorkers: maxWorkers,
		done:       make(chan struct{}),
		selector:   selector,
	}
}

func (e *Engine) Start(queue *Queue[Requester], backends []Backender) {
	// TODO: if started, should return error
	e.queue = queue
	e.backends = backends

	e.StartReportManager()
	e.StartWorkers()
}

func (e *Engine) StartWorkers() {
	for id := uint32(0); id < e.maxWorkers; id++ {
		go e.runWorker(NewWorker(id))
	}
}

func (e *Engine) runWorker(worker *Worker) {
	e.wg.Add(1)
	defer e.wg.Done()

	for {
		select {
		case <-e.done:
			return
		default:
			if request, ok := e.queue.Dequeue(); ok {
				// select backend
				backend := e.backends[e.selector.Select(request)]
				//
				e.reporter.Increment(backend.ID())
				worker.Perform(request, backend)
			} else {
				return
			}
		}
	}
}

func (e *Engine) Stop() {
	e.stopWorkers()
	e.detachQueue()
	e.stopReportManager()
	e.wg.Wait() // Shutdown gracefully
}

func (e *Engine) stopWorkers() {
	close(e.done)
}

func (e *Engine) detachQueue() {
	e.queue.Close()
}

func (e *Engine) stopReportManager() {
	e.reporter.Stop()
}

func (e *Engine) StartReportManager() {
	e.reporter = NewReportManager(5) // 5 seconds

	for _, backend := range e.backends {
		e.reporter.Register(backend.ID(), 1) // 1 minutes
	}

	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Remain items in queue: %v", e.queue.Length())
				e.reporter.Report()
				// e.ReportMetrics()
			case <-e.done:
				ticker.Stop()
				return
			}
		}
	}()
}

// func (e *Engine) ReportMetrics() {
// 	for _, backend := range e.backends {
// 		rps := backend.GetMetrics().GetRPS()
// 		log.Printf("Backend ID: %s, RPS: %f", backend.ID(), rps)
// 	}
// }
