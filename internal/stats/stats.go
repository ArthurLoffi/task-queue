package stats

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Stats struct {
	mu sync.Mutex
	Completed int
	Failed int
	InitiatedAt time.Time
}

type Snapshot struct {
	Completed int `json:"completed"`
	Failed int `json:"failed"`
	InitiatedAt string `json:"initiated_at"`
	CompletedAt string `json:"completed_at"`
}

func NewStats() *Stats {
	return &Stats{
		InitiatedAt: time.Now(),
	}
}

func (s *Stats) IncCompleted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Completed++
}

func (s *Stats) IncFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Failed++
}

func (s *Stats) Snapshot() (completed, failed int, initiatedAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Completed, s.Failed, s.InitiatedAt
}

func (s *Stats) ToJson() ([]byte, error) {
	completed, failed, initiatedAt := s.Snapshot()

	const layout = "2006-01-02 15:04:05"

	snap := Snapshot{
		Completed: completed,
		Failed: failed,
		InitiatedAt: initiatedAt.Format(layout),
		CompletedAt: time.Now().Format(layout),
	}
	return json.MarshalIndent(snap, "","  ")
}

func (s *Stats) SaveJson(path string) error {
	data, err := s.ToJson()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}