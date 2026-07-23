package entities

import (
	"task-queue/internal/dto"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	Id string `json:"id"`
	Type string `json:"task"`
	Priority int `json:"priority"`
	Status string `json:"status"`
	Attempts int `json:"attempts"`
	Payload any `json:"payload"`
	CreatedAt time.Time
}

func NewJob(req dto.CreateJobRequest) Job {
	return Job{
		Id: uuid.NewString(),
		Type: req.Type,
		Priority: req.Priority,
		Status: "queued",
		Attempts: 0,
		Payload: req.Payload,
		CreatedAt: time.Now(),
	}
}