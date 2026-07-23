package pool

import (
	"task-queue/internal/entities"
	"testing"
	"time"
)

func TestPool_Process(t *testing.T) {
	job := entities.Job{Id: "job-1"}

	start := time.Now()
	result := ProcessJob(job)
	elapsed := time.Since(start)

	if result.Id != job.Id {
		t.Errorf("Want result.Id = %s, got %s", job.Id, result.Id)
	}

	if result.Err != nil {
		t.Errorf("Expected no error, got %v", result.Err)
	}

	if result.Duration < 2*time.Second {
		t.Errorf("Expected result.Duration >= 2, got %v", result.Duration)
	}

	if elapsed < 2*time.Second {
		t.Errorf("Expected elapsed >= 2, got %v", elapsed)
	}
}