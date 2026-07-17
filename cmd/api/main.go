package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	setupRoutes(r)

	r.Run(":8080")
}