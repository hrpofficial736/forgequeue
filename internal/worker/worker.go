package worker

import (
	"context"
	"encoding/json"
	"github.com/hrpofficial736/forgequeue/internal/models"
	"github.com/hrpofficial736/forgequeue/internal/registry"
	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
	"time"
)

type WorkerInterface interface {
	StartWorker(context.Context, redis.RedisInterface)
	lookForJob(context.Context, redis.RedisInterface) *models.Job
}

type workerStruct struct {
	ri redis.RedisInterface
	rg registry.JobsRegistryInterface
}

func (w *workerStruct) lookForJobMeta(ctx context.Context) *models.Job {
	jobMetaInJsonString, err := w.ri.GetAndRemoveNextMemberFromSortedSet(ctx, "READY")

	if err != nil {
		return nil
	}

	var jobMeta *models.Job
	if err := json.Unmarshal([]byte(jobMetaInJsonString), &jobMeta); err != nil {
		return nil
	}

	return jobMeta
}

func (w *workerStruct) getJobFunction(hash string) func(...any) {
	jobFunc := w.rg.GetJob(hash)

	if jobFunc == nil {
		panic("job not found in the registry")
	}

	return jobFunc
}

func (w *workerStruct) StartWorker(ctx context.Context) {
	workerId := ""
	workerHash := "worker:" + workerId

	workerData := map[string]string{}

	w.ri.CreateHash(ctx, workerHash, workerData, 30*time.Second)

	var jobFunc func(...any)

	var jobMeta *models.Job

	for {
		jobMeta = w.lookForJobMeta(ctx)

		if jobMeta != nil {
			jobFunc = w.getJobFunction(jobMeta.ID)
			jobFunc()
		}
	}

}
