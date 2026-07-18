package main

import (
	"task-queue/cmd/controller"
	"task-queue/internal/queue"
	usecases "task-queue/internal/use-cases"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	q := queue.NewQueue(10)

	createJobUC := usecases.NewCreateJobUseCase(q)
	ctrl := controller.NewController(createJobUC)

	r.POST("/job", ctrl.CreateJob)

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
		})
	})

	r.Run(":8080")
}