package services

import (
	"videoEncoder/app/repositories"
	"videoEncoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) changeJobStatus(status string) error {
	var err error

	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobService) failJob(error error) error {
	j.Job.Status = "FAILED"
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return error
}