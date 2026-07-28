package main

import (
	"log"
	pool "task-queue/internal/worker-pool"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	agingCfg := pool.AgingConfig{
		Interval: 10 * time.Second,
		Deadlines: map[int]time.Duration{
			1: 5 * time.Minute,
			2: 2 * time.Minute,
		},
		MaxPriority: 3,
	}

	r := gin.Default()
	p := pool.NewPool(3, 10, &agingCfg)
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