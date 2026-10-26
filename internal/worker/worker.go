package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hrpofficial736/forgequeue/internal/models"
	"github.com/hrpofficial736/forgequeue/internal/registry"
	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
	"time"
)

type WorkerInterface interface {
	StartWorker(context.Context)
	lookForJobMeta(context.Context) *models.Job
	getJobFunction(string) (*models.JobPayload, error)
	createOrUpdateWorker(context.Context, string, *models.Worker) error
}

type workerStruct struct {
	ri redis.RedisInterface
	rg registry.JobsRegistryInterface
}

func NewWorkerInterface(ri redis.RedisInterface, rg registry.JobsRegistryInterface) WorkerInterface {
	return &workerStruct{ri, rg}
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

func (w *workerStruct) getJobFunction(hash string) (*models.JobPayload, error) {
	jobFunc, err := w.rg.GetJob(hash)

	if err != nil {
		return &models.JobPayload{}, fmt.Errorf("no function for given job ID")
	}

	return jobFunc, nil
}

func (w *workerStruct) createOrUpdateWorker(ctx context.Context, hash string, data *models.Worker) error {
	workerJson, err := json.Marshal(data)

	if err != nil {
		fmt.Println(fmt.Errorf("error converting worker's data to json: %s", err))
		return err
	}

	// creating worker hash in redis
	err = w.ri.SetKeyValue(ctx, hash, workerJson, data.HeartbeatInterval)

	if err != nil {
		fmt.Println(fmt.Errorf("error creating worker in redis: %s", err))
		return err
	}

	return nil
}

func (w *workerStruct) StartWorker(ctx context.Context) {
	// creating worker hash
	workerHash := fmt.Sprintf("worker:%s", uuid.New().String())

	// creating worker model
	workerData := &models.Worker{
		ID:                workerHash,
		HeartbeatInterval: 30 * time.Second,
	}

	err := w.createOrUpdateWorker(ctx, workerHash, workerData)

	time.AfterFunc(30*time.Second, func() {
		err := w.createOrUpdateWorker(ctx, workerHash, workerData)
		if err != nil {
			return
		}
	})

	if err != nil {
		return
	}

	job := &models.JobPayload{}

	for {
		jobMeta := w.lookForJobMeta(ctx)

		if jobMeta != nil {
			job, err = w.getJobFunction(jobMeta.ID)
			if err != nil {
				fmt.Println(err)
				return
			}

			workerData.CurrentJobID = jobMeta.ID

			err := w.createOrUpdateWorker(ctx, workerHash, workerData)
			if err != nil {
				return
			}

			job.Handler(job.Payload)
		}
	}

}
