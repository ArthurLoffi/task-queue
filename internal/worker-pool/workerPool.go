package pool

import (
	"context"
	"log"
	"sync"
	"task-queue/internal/entities"
	"task-queue/internal/stats"
	"time"
)

type Pool struct {
	Jobs chan entities.Job
	Results chan Result
	wg sync.WaitGroup
	s *stats.Stats
}

func NewPool(numWorkers int, bufferSize int) *Pool {
	p := &Pool{
		Jobs: make(chan entities.Job, bufferSize),
		Results: make(chan Result, bufferSize),
		s: stats.NewStats(),
	}
	for i := range numWorkers {
		p.wg.Add(1)
		go p.worker(i)
	}
	
	return p
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for j := range p.Jobs {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			3 * time.Second)
		
		result := p.processWithContext(ctx, id, j)
		cancel() // Libera recursos do context
		p.Results <- result
	}
}

func (p *Pool) processWithContext(ctx context.Context, workerID int, j entities.Job) Result {
	// Ele ainda está rodando no background
	// Passar o context para o processJob e resolver internamente
	done := make(chan Result, 1)

	go func() {
		done <- ProcessJob(j)
	}()

	select {
		case result := <- done:
			p.s.IncCompleted()
			return result
		case <- ctx.Done():
			log.Printf("Timeout worker %d in job %s", workerID, j.Id)
			p.s.IncFailed()
			return Result{
				Id: j.Id,
				Err: ctx.Err(),
			}
	}
}

func (p *Pool) Submit(j entities.Job) {
	p.Jobs <- j
}

func (p *Pool) Shutdown() {
	close(p.Jobs)
	p.wg.Wait()
	close(p.Results)
}

func (p *Pool) Stats() *stats.Stats {
	return p.s
}