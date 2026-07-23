package dto

type CreateJobRequest struct {
	Type string `json:"type"`
	Priority int `json:"priority"`
	Status string `json:"status"`
	Payload any `json:"payload"`
}