package main

import (
	"context"
	"fmt"
	"github.com/hrpofficial736/forgequeue/internal/models"
	"github.com/hrpofficial736/forgequeue/internal/recoverer"
	"github.com/hrpofficial736/forgequeue/internal/registry"
	"github.com/hrpofficial736/forgequeue/internal/scheduler"
	redis "github.com/hrpofficial736/forgequeue/internal/storage/redis"
	"github.com/hrpofficial736/forgequeue/internal/worker"
	"github.com/hrpofficial736/forgequeue/lib"
)

func main() {
	ctx := context.Background()

	client := redis.LoadRedis()

	ri := redis.NewRedisInterface(client)
	rg := registry.NewJobsRegistryInterface()

	w := worker.NewWorkerInterface(ri, rg)
	s := scheduler.NewSchedulerInterface(ri)
	r := recoverer.NewRecovererInterface(ri)

	engineInterface := lib.NewEngineInterface(w, s, r)

	go engineInterface.StartEngine(ctx)

	// sample job
	jc := models.JobConfiguration{DefaultPriority: 100, DefaultMaxRetries: 5}

	enqJobInterface := lib.NewEnqueueJobInterface(ri, rg, jc)

	enqJobInterface.EnqueueInstantJob(ctx, &models.JobPayload{
		Name: "print_hello",
		Payload: map[string]interface{}{
			"name": "Harry",
		},
		Handler: func(m map[string]any) {
			fmt.Println("hello", m["name"])
		},
	})

	select {}

}
