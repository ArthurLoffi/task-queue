package controller

import (
	"log"
	"net/http"
	"task-queue/internal/dto"
	e "task-queue/internal/entities"
	usecases "task-queue/internal/use-cases"

	"github.com/gin-gonic/gin"
)

func CreateJobController(c *gin.Context) {
	var body dto.CreateJobRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	newJob := e.Job{
		Type: body.Type,
		Priority: body.Priority,
		Status: body.Status,
		Payload: body.Payload,
	}

	err := usecases.CreateJob(newJob)
	if err != nil {
		log.Printf("Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
}