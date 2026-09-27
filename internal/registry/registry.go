package registry

import "fmt"

var jobsMappings = make(map[string]func(...any))

type JobsRegistryInterface interface {
	GetJob(string) func(...any)
	CreateJob(string, func(...any)) error
	DeleteJob(string) error
}

func NewJobsRegistryInterface() JobsRegistryInterface {
	return &jobRegistry{}
}

type jobRegistry struct{}

func (j *jobRegistry) GetJob(hash string) func(...any) {
	if job, exists := jobsMappings[hash]; exists {
		return job
	}

	return nil
}

func (j *jobRegistry) CreateJob(hash string, job func(...any)) error {
	if _, exists := jobsMappings[hash]; exists {
		return fmt.Errorf("job already exists")
	}

	jobsMappings[hash] = job
	return nil
}

func (j *jobRegistry) DeleteJob(hash string) error {
	if _, exists := jobsMappings[hash]; !exists {
		return fmt.Errorf("job not found")
	}

	delete(jobsMappings, hash)

	return nil
}
