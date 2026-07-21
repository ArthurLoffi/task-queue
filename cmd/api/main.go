package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-queue/cmd/controller"
	usecases "task-queue/internal/use-cases"
	pool "task-queue/internal/worker-pool"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// q := queue.NewQueue(10)

	p := pool.NewPool(3, 10)
	s := p.Stats()

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

	srv := http.Server{
		Addr: ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error to run server: %v", err)
		}
	}()


	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error to shutdown the server: %v", err)
	}

	p.Shutdown()

	if err := s.SaveJson("db.json"); err != nil {
		log.Printf("Error to save stats in json: %v", err)
	}

	log.Println("Server shutdown successfully!")
}