package pool

import (
	"task-queue/internal/entities"
	"time"
)

type Result struct {
	Id string
	Priority int
	Err error
	Duration time.Duration
}

func ProcessJob(j entities.Job) Result {
	start := time.Now()
	time.Sleep(2 * time.Second) // Simulando que ta fazendo algo pesado
	return Result{
		Id: j.Id,
		Priority: j.Priority,
		Duration: time.Since(start),
	}
}