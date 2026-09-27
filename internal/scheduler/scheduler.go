package scheduler

import (
	"context"
	"time"

	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
)

type SchedulerInterface interface {
	scheduleCronJobs(ctx context.Context) error
	scheduleDelayedJobs(ctx context.Context) error
	ScheduleJobs(ctx context.Context)
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
	cronJobs, err := s.ri.GetMembersWithASpecificScore(
		ctx,
		"SCHEDULES",
		float64(currTime),
		float64(currTime),
		float64(currTime),
	)

	if err != nil {
		return err
	}

	for _, job := range cronJobs {
		err := s.ri.AddToSortedSet(ctx, "READY", float64(currTime), job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scheduler) scheduleDelayedJobs(ctx context.Context) error {
	// taking jobs out of the DELAYED queue

	currTime := time.Now().UnixMilli()
	delayedJobs, err := s.ri.GetMembersWithASpecificScore(
		ctx,
		"DELAYED",
		float64(currTime),
		0,
		float64(currTime),
	)

	if err != nil {
		return err
	}

	for _, job := range delayedJobs {
		err := s.ri.AddToSortedSet(ctx, "READY", float64(currTime), job)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *scheduler) ScheduleJobs(ctx context.Context) {
	for {
		go s.scheduleCronJobs(ctx)
		go s.scheduleDelayedJobs(ctx)
	}
}
