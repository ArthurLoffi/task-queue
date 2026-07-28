package pool

import (
	"sync"
	"task-queue/internal/entities"
	"testing"
	"time"
)

var agingCfg = AgingConfig{
		Interval: 10 * time.Second,
		Deadlines: map[int]time.Duration{
			1: 5 * time.Minute,
			2: 2 * time.Minute,
		},
		MaxPriority: 3,
	}

// fakeProcessor permite controlar, por teste, quanto tempo o
// processamento leva e se ele retorna erro
type fakeProcessor struct {
	mu    sync.Mutex
	delay time.Duration
	err   error
}

func (f *fakeProcessor) Process(j entities.Job) Result {
	f.mu.Lock()
	delay := f.delay
	err := f.err
	f.mu.Unlock()

	if delay > 0 {
		time.Sleep(delay)
	}

	return Result{
		Id:  j.Id,
		Err: err,
	}
}

// Teste do processador default
// Func para testes
func TestDefaultProcessor_Process(t *testing.T) {
	result := defaultProcessor{}.Process(entities.Job{Id: "job-1"})

	if result.Id != "job-1" {
		t.Errorf("Want result.Id = job-1, got %s", result.Id)
	}
	if result.Err != nil {
		t.Errorf("Expected no error, got: %v", result.Err)
	}
}

// Valida essa Func para garantir
// que todos os campos sejam preenchidos
func TestPool_NewPool(t *testing.T) {
	pool := NewPool(3, 10, &agingCfg)
	defer pool.Shutdown()

	if pool.Jobs == nil || pool.processor == nil || pool.Results == nil || pool.s == nil {
		t.Errorf("Expected a completed pool struct")
	}
}

// Testa se está funcionando o envio e o processamento de um job
// Verifica se o resultado retornado é o mesmo do id do job criado
// Também verifica retorno de erro e timeout
func TestPool_SubmitAndProcess(t *testing.T) {
	p := NewPoolWithProcessor(3, 5, &fakeProcessor{delay: 50 * time.Millisecond}, &agingCfg)

	job := entities.Job{Id: "job-1"}
	p.Submit(job)

	select {
	case result := <-p.Results:
		if result.Id != job.Id {
			t.Errorf("Want result.Id = %s, got %s", job.Id, result.Id)
		}
		if result.Err != nil {
			t.Fatalf("Expected without error, got: %v", result.Err)
		}
	case <-time.After(1 * time.Second):
		t.Errorf("Timeout waiting job result")
	}

	p.Shutdown()

	completed, failed, _ := p.Stats().Snapshot()
	if completed != 1 {
		t.Errorf("Expected one job completed, got %d", completed)
	}
	if failed != 0 {
		t.Errorf("Expected zero job failed, got %d", failed)
	}
}

// Verifica se todos os jobs são processados corretamente
// por workers diferentes
func TestPool_MultipleJobs(t *testing.T) {
	p := NewPoolWithProcessor(3, 10, &fakeProcessor{delay: 50 * time.Millisecond}, &agingCfg)

	const total = 10
	for i := range total {
		p.Submit(entities.Job{Id: string(rune('a' + i))})
	}

	received := 0
	timeout := time.After(3 * time.Second)

loop:
	for received < total {
		select {
		case <-p.Results:
			received++
		case <-timeout:
			break loop
		}
	}

	p.Shutdown()

	if received != total {
		t.Errorf("Expected %d results, got %d", total, received)
	}
}

// Garante que após um p.Shutdown, o canal
// Results é fechado sem ficar com goroutines travados
func TestPool_ShutdownClosesResults(t *testing.T) {
	p := NewPoolWithProcessor(1, 1, &fakeProcessor{}, &agingCfg)

	p.Submit(entities.Job{Id: "job-x"})
	<-p.Results

	p.Shutdown()

	select {
	case _, ok := <-p.Results:
		if ok {
			t.Errorf("Expected channel Results closed, but it still return a value")
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting the channel Results close")
	}
}

// Falta adicionar uma injeção para criar um job
// lento que simule o timeout
func TestPool_Timeout(t *testing.T) {
	p := NewPoolWithProcessor(1, 1, &fakeProcessor{delay: 4 * time.Second}, &agingCfg)

	p.Submit(entities.Job{Id: "job-slow"})

	select {
	case result := <-p.Results:
		if result.Err == nil {
			t.Errorf("Expected timeout error, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout waiting job result")
	}

	p.Shutdown()

	_, failed, _ := p.Stats().Snapshot()
	if failed != 1 {
		t.Errorf("Expected one job failed, got %d", failed)
	}
}

// Test para verificar se está aumenta a priority
// para previnir starvation, não está 100%
// pois preciso testar o log também dessa func
func TestPool_StartAgingPromotesExpiredJobs(t *testing.T) {
	cfg := &AgingConfig{
		Interval: 30 * time.Millisecond,
		Deadlines: map[int]time.Duration{1: 50 * time.Millisecond},
		MaxPriority: 3,
	}

	proc := &fakeProcessor{}
	p := NewPoolWithProcessor(1, 10, proc, cfg)
	defer p.Shutdown()

	// Só ocupa para testar o aging
	proc.mu.Lock()
	proc.delay = 200 * time.Millisecond
	proc.mu.Unlock()

	p.Submit(entities.Job{Id: "blocker", Priority: 3, CreatedAt: time.Now()})

	time.Sleep(20 * time.Millisecond)

	proc.mu.Lock()
	proc.delay = 0
	proc.mu.Unlock()

	jobB := entities.Job{Id: "B", Priority: 1, CreatedAt: time.Now()}
	jobA := entities.Job{Id: "A", Priority: 1, CreatedAt: time.Now().Add(-100 * time.Millisecond)} // já expirado

	p.Submit(jobB)
	p.Submit(jobA)

	var order []string
	for i := 0; i < 3; i++ {
		select {
		case r := <-p.Results:
			order = append(order, r.Id)
		case <-time.After(2 * time.Second):
			t.Fatalf("Timeout waiting results, partial order: %v", order)
		}
	}

	if order[0] != "blocker" {
		t.Fatalf("Expected blocker first, got %v", order)
	}
	if order[1] != "A" {
		t.Errorf("Expected job A (promoted by aging) before the job B, got order: %v", order)
	}
}