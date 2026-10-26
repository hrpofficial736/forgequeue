package models

import "time"

type JobPayload struct {
	Name    string
	Payload map[string]interface{}
	Handler func(map[string]interface{})
}

type JobStatus int

const (
	READY JobStatus = iota
	ACTIVE
	SUCCEEDED
	FAILED
	DELAYED
	SCHEDULED
)

type Job struct {
	ID           string    `redis:"id"`
	Name         string    `redis:"name"`
	QueueName    string    `redis:"queueName"`
	Status       JobStatus `redis:"status"`
	LeaseOwner   string    `redis:"leaseOwner"`
	LeaseExpiry  time.Time `redis:"leaseExpiry"`
	Priority     int       `redis:"priority"`
	RetriesCount int       `redis:"retriesCount"`
	MaxRetries   int       `redis:"maxRetries"`

	// if job is DELAYED
	After time.Duration `redis:"after"`

	// if job is SCHEDULED
	At time.Time `redis:"at"`

	CreatedAt time.Time `redis:"createdAt"`
}

type JobConfiguration struct {
	DefaultPriority   int
	DefaultMaxRetries int
}

type Worker struct {
	ID                string        `redis:"id"`
	HeartbeatInterval time.Duration `redis:"heartBeatInterval"`
	CurrentJobID      string        `redis:"currentJobId"`
}
