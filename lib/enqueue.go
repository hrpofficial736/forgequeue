package lib

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hrpofficial736/forgequeue/internal/models"
	"github.com/hrpofficial736/forgequeue/internal/registry"
	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
)

type EnqueueJobInterface interface {
	EnqueueInstantJob(context.Context, *models.JobPayload) error
	EnqueueDelayedJob(context.Context, *models.JobPayload, time.Duration) error
	EnqueueCronJob(context.Context, *models.JobPayload, time.Time) error
}

type enqueueJobStruct struct {
	ri redis.RedisInterface
	rg registry.JobsRegistryInterface
	jc models.JobConfiguration
}

func NewEnqueueJobInterface(
	ri redis.RedisInterface,
	rg registry.JobsRegistryInterface,
	jc models.JobConfiguration,
) EnqueueJobInterface {
	return &enqueueJobStruct{ri, rg, jc}
}

func (e *enqueueJobStruct) EnqueueInstantJob(ctx context.Context, job *models.JobPayload) error {

	jobMetaData := models.Job{
		ID:           "job-" + uuid.NewString(),
		Name:         job.Name,
		QueueName:    "READY",
		Status:       models.READY,
		Priority:     e.jc.DefaultPriority,
		RetriesCount: 0,
		MaxRetries:   e.jc.DefaultMaxRetries,
		CreatedAt:    time.Now(),
	}

	err := e.rg.CreateJob(jobMetaData.ID, job)
	if err != nil {
		return err
	}

	jobMetaDataInBytes, err := json.Marshal(jobMetaData)

	if err != nil {
		return err
	}

	err = e.ri.AddMemberToSortedSet(ctx, "READY", float64(e.jc.DefaultPriority), string(jobMetaDataInBytes))

	if err != nil {
		return err
	}

	return nil

}

func (e *enqueueJobStruct) EnqueueDelayedJob(ctx context.Context, job *models.JobPayload, after time.Duration) error {
	jobMetaData := models.Job{
		ID:           "job-" + uuid.NewString(),
		Name:         job.Name,
		QueueName:    "DELAYED",
		Status:       models.DELAYED,
		Priority:     e.jc.DefaultPriority,
		RetriesCount: 0,
		MaxRetries:   e.jc.DefaultMaxRetries,
		After:        after,
		CreatedAt:    time.Now(),
	}

	err := e.rg.CreateJob(jobMetaData.ID, job)
	if err != nil {
		return err
	}

	jobMetaDataInBytes, err := json.Marshal(jobMetaData)

	if err != nil {
		return err
	}

	err = e.ri.AddMemberToSortedSet(ctx, "DELAYED", float64(e.jc.DefaultPriority), string(jobMetaDataInBytes))

	if err != nil {
		return err
	}

	return nil

}

func (e *enqueueJobStruct) EnqueueCronJob(ctx context.Context, job *models.JobPayload, at time.Time) error {
	jobMetaData := models.Job{
		ID:           "job-" + uuid.NewString(),
		Name:         job.Name,
		QueueName:    "CRON",
		Status:       models.SCHEDULED,
		Priority:     e.jc.DefaultPriority,
		RetriesCount: 0,
		MaxRetries:   e.jc.DefaultMaxRetries,
		At:           at,
		CreatedAt:    time.Now(),
	}

	err := e.rg.CreateJob(jobMetaData.ID, job)
	if err != nil {
		return err
	}

	jobMetaDataInBytes, err := json.Marshal(jobMetaData)

	if err != nil {
		return err
	}

	err = e.ri.AddMemberToSortedSet(ctx, "CRON", float64(e.jc.DefaultPriority), string(jobMetaDataInBytes))

	if err != nil {
		return err
	}

	return nil

}
