package usecases

import (
	"task-queue/internal/dto"
	"task-queue/internal/entities"
	pool "task-queue/internal/worker-pool"
)

type CreateJobUseCase struct {
	pool *pool.Pool
}

func NewCreateJobUseCase(p *pool.Pool) *CreateJobUseCase {
	return &CreateJobUseCase{
		pool: p,
	}
}

func (uc *CreateJobUseCase) Execute(req dto.CreateJobRequest) (entities.Job, error) {
	job := entities.NewJob(req)
	uc.pool.Submit(job)
	return job, nil
}