package registry

import (
	"fmt"

	"github.com/hrpofficial736/forgequeue/internal/models"
)

var jobsMappings *map[string]*models.JobPayload = &map[string]*models.JobPayload{}

type JobsRegistryInterface interface {
	GetJob(string) (*models.JobPayload, error)
	CreateJob(string, *models.JobPayload) error
	DeleteJob(string) error
}

func NewJobsRegistryInterface() JobsRegistryInterface {
	return &jobRegistry{}
}

type jobRegistry struct{}

func (j *jobRegistry) GetJob(hash string) (*models.JobPayload, error) {
	if job, exists := (*jobsMappings)[hash]; exists {
		return job, nil
	}

	return &models.JobPayload{}, fmt.Errorf("no function found for the given job ID")
}

func (j *jobRegistry) CreateJob(hash string, job *models.JobPayload) error {
	if _, exists := (*jobsMappings)[hash]; exists {
		return fmt.Errorf("job already exists")
	}

	(*jobsMappings)[hash] = job
	return nil
}

func (j *jobRegistry) DeleteJob(hash string) error {
	if _, exists := (*jobsMappings)[hash]; !exists {
		return fmt.Errorf("job not found")
	}

	delete((*jobsMappings), hash)

	return nil
}
