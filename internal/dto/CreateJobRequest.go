package dto

import "task-queue/internal/entities"

type CreateJobRequest struct {
	Type string `json:"type"`
	Priority string `json:"priority"`
	Status string `json:"status"`
	Payload entities.Payload `json:"payload"`
}