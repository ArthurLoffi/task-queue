package controller

import (
	"net/http"
	"task-queue/internal/dto"
	usecases "task-queue/internal/use-cases"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	createJobUC *usecases.CreateJobUseCase
}

func NewController(createJobUC *usecases.CreateJobUseCase) *Controller {
	return &Controller{
		createJobUC: createJobUC,
	}
}

func (c *Controller) CreateJob(ctx *gin.Context) {
	var req dto.CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	job, err := c.createJobUC.Execute(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, job)
}