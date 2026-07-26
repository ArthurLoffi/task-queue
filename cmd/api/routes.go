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

func SetupRoutes(r *gin.Engine, p *pool.Pool) {
    createJobUC := usecases.NewCreateJobUseCase(p)
    ctrl := controller.NewController(createJobUC)

    r.POST("/job", ctrl.CreateJob)
    r.GET("/healthy", func(c *gin.Context) {
        c.JSON(200, gin.H{"success": true})
    })
}

func Run(r *gin.Engine, p *pool.Pool) {
    srv := &http.Server{Addr: ":8080", Handler: r}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Error to run server: %v", err)
        }
    }()

    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh

    log.Println("Shutting down the server")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Printf("Error to shutdown the server: %v", err)
    }

    p.Shutdown()
}