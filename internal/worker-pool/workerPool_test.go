package pool

import (
	"sync"
	"task-queue/internal/entities"
	"testing"
	"time"
)

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

// Testa se está funcionando o envio e o processamento de um job
// Verifica se o resultado retornado é o mesmo do id do job criado
// Também verifica retorno de erro e timeout
func TestPool_SubmitAndProcess(t *testing.T) {
	p := NewPoolWithProcessor(3, 5, &fakeProcessor{})

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
	case <-time.After(2 * time.Second):
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
	p := NewPoolWithProcessor(3, 10, &fakeProcessor{})

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
	p := NewPoolWithProcessor(1, 1, &fakeProcessor{})

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
func TestPool_TImeout(t *testing.T) {
	t.Skip("Tem que fazer a simulação de um job lento em especifico")
}
