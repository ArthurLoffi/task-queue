package main

import (
	"log"
	"task-queue/cmd/controller"
	usecases "task-queue/internal/use-cases"
	pool "task-queue/internal/worker-pool"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// q := queue.NewQueue(10)

	p := pool.NewPool(3, 10)

	createJobUC := usecases.NewCreateJobUseCase(p)
	ctrl := controller.NewController(createJobUC)

	go func()  {
		for result := range p.Results {
			log.Printf("Job ID: %s | Processed in: %s", result.Id, result.Duration)
		}
	}()

	r.POST("/job", ctrl.CreateJob)

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
		})
	})

	r.Run(":8080")
}