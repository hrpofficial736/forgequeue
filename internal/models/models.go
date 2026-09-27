package models

import "time"

type JobStatus string

const (
	READY     JobStatus = "READY"
	ACTIVE    JobStatus = "ACTIVE"
	SUCCEEDED JobStatus = "SUCCEEDED"
	FAILED    JobStatus = "FAILED"
	DELAYED   JobStatus = "DELAYED"
)

type Job struct {
	ID           string
	Name         string
	HandlerName  string
	QueueName    string
	Payload      any
	Status       JobStatus
	LeaseOwner   string
	LeaseExpiry  time.Time
	Priority     int
	RetriesCount int
	MaxRetries   int
	Queue        string
	CreatedAt    time.Time
}
