package pool

import (
	"context"
	"log"
	"sync"
	"task-queue/internal/entities"
	"task-queue/internal/queue"
	"task-queue/internal/stats"
	"time"
)

type Processor interface {
	Process(j entities.Job) Result
}

type defaultProcessor struct{}

func (defaultProcessor) Process(j entities.Job) Result {
	return ProcessJob(j)
}

type Pool struct {
	Jobs *queue.PriorityQueue
	Results chan Result
	wg sync.WaitGroup
	s *stats.Stats
	processor Processor
	agingCfg *AgingConfig
	stopAging chan struct{}
}

type AgingConfig struct {
	Interval time.Duration
	Deadlines map[int]time.Duration
	MaxPriority int
}

// Novo metodo de criar pool com uma abstração
// ao inves de usar diretamente o processJob.
// Serve mais para test
func NewPool(numWorkers int, bufferSize int, agingCfg *AgingConfig) *Pool {
	return NewPoolWithProcessor(numWorkers, bufferSize, defaultProcessor{}, agingCfg)
}

func NewPoolWithProcessor(numWorkers int, bufferSize int, processor Processor, agingCfg *AgingConfig) *Pool {
	p := &Pool{
		Jobs: queue.NewPriorityQueue(),
		Results: make(chan Result, bufferSize),
		s: stats.NewStats(),
		processor: processor,
		agingCfg: agingCfg,
		stopAging: make(chan struct{}),
	}

	for i := range numWorkers {
		p.wg.Add(1)
		go p.worker(i)
	}

	if agingCfg != nil {
		go p.startAging()
	} 
	
	return p
}

func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		j, ok := p.Jobs.Pop()

		if !ok {
			return
		}

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
		done <- p.processor.Process(j)
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
	p.Jobs.Push(j)
}

func (p *Pool) startAging() {
	ticket := time.NewTicker(p.agingCfg.Interval)
	defer ticket.Stop()

	for {
		select{
		case <-p.stopAging:
			return
		case <-ticket.C:
			n := p.Jobs.PromoteExpired(p.agingCfg.Deadlines, p.agingCfg.MaxPriority)
			if n >0 {
				log.Printf("aging: %d jobs promoted by deadline", n)
			}
		}
	}
}

func (p *Pool) Shutdown() {
	p.Jobs.Close()
	p.wg.Wait()
	close(p.Results)
}

func (p *Pool) Stats() *stats.Stats {
	return p.s
}