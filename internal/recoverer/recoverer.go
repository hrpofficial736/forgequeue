package recoverer

import (
	"context"
	"time"

	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
)

type RecovererInterface interface {
	RecoverJobs(context.Context) error
}

type recoverer struct {
	ri redis.RedisInterface
}

func NewRecovererInterface(ri redis.RedisInterface) RecovererInterface {
	return &recoverer{ri}
}

func (r *recoverer) RecoverJobs(ctx context.Context) error {
	// calculating the current time
	currTime := time.Now().UnixMilli()

	// getting all jobs that have expired lease(currTime == lease time)
	jobs, err := r.ri.GetMembersWithASpecificScore(
		ctx,
		"ACTIVE",
		float64(currTime),
		0,
		float64(currTime),
	)

	if err != nil {
		return err
	}

	// removing all those jobs from ACTIVE queue
	for _, job := range jobs {
		err := r.ri.RemoveASpecificMemberFromSortedSet(ctx, "ACTIVE", job)
		if err != nil {
			return err
		}
	}

	// pushing jobs back to READY queue
	for _, job := range jobs {
		err := r.ri.AddToSortedSet(ctx, "READY", float64(currTime), job)
		if err != nil {
			return err
		}
	}

	return nil
}
