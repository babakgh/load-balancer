package engine

import (
	"errors"
	"log"
)

var ErrWorkerFailed = errors.New("ErrWorkerFailed")

type Worker struct {
	ID uint32
}

func NewWorker(id uint32) *Worker {
	return &Worker{
		ID: id,
	}
}

func (w *Worker) Perform(request Requester, backend Backender) error {
	if err := backend.Process(request); err != nil {
		log.Printf("Error processing request by worker %v: %v", request.ID(), err)
		return ErrWorkerFailed
	}

	return nil
}
