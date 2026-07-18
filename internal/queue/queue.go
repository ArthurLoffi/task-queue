package queue

import "task-queue/internal/entities"

type Queue struct {
	jobs chan entities.Job
}

func NewQueue(bufferSize int) *Queue {
	return &Queue{
		jobs: make(chan entities.Job, bufferSize),
	}
}

func (q *Queue) Push(job entities.Job) {
	q.jobs <- job
}

func (q *Queue) Jobs() <- chan entities.Job {
	return q.jobs
}