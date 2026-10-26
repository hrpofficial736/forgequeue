package lib

import (
	"context"
	"fmt"

	"github.com/hrpofficial736/forgequeue/internal/recoverer"
	"github.com/hrpofficial736/forgequeue/internal/scheduler"
	"github.com/hrpofficial736/forgequeue/internal/worker"
)

type EngineInterface interface {
	StartEngine(context.Context)
}

type engine struct {
	w worker.WorkerInterface
	s scheduler.SchedulerInterface
	r recoverer.RecovererInterface
}

func NewEngineInterface(
	w worker.WorkerInterface,
	s scheduler.SchedulerInterface,
	r recoverer.RecovererInterface,
) EngineInterface {
	return &engine{w, s, r}
}

func (e *engine) StartEngine(ctx context.Context) {
	// starting a worker
	e.w.StartWorker(ctx)

	fmt.Println("worker running...")
	// starting a scheduler
	e.s.StartScheduler(ctx)

	fmt.Println("scheduler running...")
	// starting a recoverer
	e.r.StartRecoverer(ctx)

	fmt.Println("recoverer running...")
}
