package main

import (
	"log"
	pool "task-queue/internal/worker-pool"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	p := pool.NewPool(3, 10)
	s := p.Stats()

	go func()  {
		for result := range p.Results {
			log.Printf("Job ID: %s | Priority: %d | Processed in: %s", result.Id, result.Priority, result.Duration)
		}
	}()

	SetupRoutes(r, p)
	Run(r, p)

	if err := s.SaveJson("db.json"); err != nil {
		log.Printf("Error to save stats in json: %v", err)
	}

	log.Println("Server shutdown successfully!")
}