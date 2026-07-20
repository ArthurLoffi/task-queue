package pool

import (
	"sync"
	"task-queue/internal/entities"
)

type Pool struct {
	Jobs chan entities.Job
	Results chan Result
	wg sync.WaitGroup
}

func NewPool(numWorkers int, bufferSize int) *Pool {
	p := &Pool{
		Jobs: make(chan entities.Job, bufferSize),
		Results: make(chan Result, bufferSize),
	}
	for range numWorkers {
		p.wg.Add(1)
		go func ()  {
			defer p.wg.Done()
			for j := range p.Jobs {
				p.Results <- ProcessJob(j)
			}
		}()
	}
	return p
}

func (p *Pool) Submit(j entities.Job) {
	p.Jobs <- j
}

func (p *Pool) Shutdown() {
	close(p.Jobs)
	p.wg.Wait()
	close(p.Results)
}