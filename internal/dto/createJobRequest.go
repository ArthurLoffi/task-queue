package dto

type CreateJobRequest struct {
	Type string `json:"type"`
	Priority string `json:"priority"`
	Status string `json:"status"`
	Payload any `json:"payload"`
}