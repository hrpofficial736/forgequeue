package recoverer

import (
	"context"
	"encoding/json"
	"github.com/hrpofficial736/forgequeue/internal/models"
	"github.com/hrpofficial736/forgequeue/internal/storage/redis"
	"time"
)

type RecovererInterface interface {
	StartRecoverer(context.Context) error
}

type recoverer struct {
	ri redis.RedisInterface
}

func NewRecovererInterface(ri redis.RedisInterface) RecovererInterface {
	return &recoverer{ri}
}

func (r *recoverer) StartRecoverer(ctx context.Context) error {
	// calculating the current time
	currTime := time.Now().UnixMilli()

	// getting all jobs that have expired lease(currTime == lease time)
	jobs, err := r.ri.GetARangeOfMembers(
		ctx,
		"ACTIVE",
		0,
		currTime,
	)

	if err != nil {
		return err
	}

	for _, jobAsString := range jobs {
		var job models.Job

		if err := json.Unmarshal([]byte(jobAsString), &job); err != nil {
			return err
		}

		// removing all those jobs from ACTIVE queue
		err := r.ri.RemoveASpecificMemberFromSortedSet(ctx, "ACTIVE", job.ID)
		if err != nil {
			return err
		}

		jobInBytes, err := json.Marshal(job)

		if err != nil {
			return err
		}

		// pushing jobs back to READY queue
		err = r.ri.AddMemberToSortedSet(ctx, "READY", float64(currTime), string(jobInBytes))
		if err != nil {
			return err
		}
	}

	return nil
}
