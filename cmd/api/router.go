package main

import (
	"task-queue/cmd/controller"

	"github.com/gin-gonic/gin"
)

func setupRoutes(r *gin.Engine)  {
	r.POST("/job", controller.CreateJobController)

	r.GET("/healthy", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
		})
	})
}