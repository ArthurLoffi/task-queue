package usecases

import (
	"task-queue/internal/dto"
	"task-queue/internal/entities"
	"task-queue/internal/queue"

)

type CreateJobUseCase struct {
	queue *queue.Queue
}

func NewCreateJobUseCase(q *queue.Queue) *CreateJobUseCase {
	return &CreateJobUseCase{
		queue: q,
	}
}

func (uc *CreateJobUseCase) Execute(req dto.CreateJobRequest) (entities.Job, error) {
	job := entities.NewJob(req)
	uc.queue.Push(job)
	return job, nil
}