package scheduler

import (
	"context"
	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
	"time"
)

type SchedulerInterface interface {
	scheduleCronJobs(ctx context.Context) error
	scheduleDelayedJobs(ctx context.Context) error
	StartScheduler(ctx context.Context)
}

type scheduler struct {
	ri redis.RedisInterface
}

func NewSchedulerInterface(ri redis.RedisInterface) SchedulerInterface {
	return &scheduler{ri}
}

func (s *scheduler) scheduleCronJobs(ctx context.Context) error {
	// taking jobs out of the SCHEDULES queue
	currTime := time.Now().UnixMilli()
	cronJobs, err := s.ri.GetARangeOfMembers(
		ctx,
		"SCHEDULES",
		currTime,
		currTime,
	)

	if err != nil {
		return err
	}

	for _, job := range cronJobs {
		err := s.ri.AddMemberToSortedSet(ctx, "READY", float64(currTime), job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scheduler) scheduleDelayedJobs(ctx context.Context) error {
	// taking jobs out of the DELAYED queue

	currTime := time.Now().UnixMilli()
	delayedJobs, err := s.ri.GetARangeOfMembers(
		ctx,
		"DELAYED",
		currTime,
		currTime,
	)

	if err != nil {
		return err
	}

	for _, job := range delayedJobs {
		err := s.ri.AddMemberToSortedSet(ctx, "READY", float64(currTime), job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scheduler) StartScheduler(ctx context.Context) {
	for {
		go s.scheduleCronJobs(ctx)
		go s.scheduleDelayedJobs(ctx)
	}
}
