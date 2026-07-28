package queue

import (
	"container/heap"
	"sync"
	"task-queue/internal/entities"
	"time"
)

type item struct {
	job entities.Job
	index int
	promoted bool
}

type itemHeap []*item

func (h itemHeap) Len() int {
	return len(h)
}

func (h itemHeap) Less(i, j int) bool {
	return h[i].job.Priority > h[j].job.Priority
}

func (h itemHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *itemHeap) Push(x any) {
	it := x.(*item)
	it.index = len(*h)
	*h = append(*h, it)
}

func (h *itemHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	old[n-1] = nil
	it.index = -1
	*h = old[:n-1]
	return it
}

// PriorityQueue é uma fila de prioridade thread-safe. Múltiplos
// workers podem chamar Pop() concorrentemente, e múltiplos produtores
// podem chamar Push() concorrentemente.
type PriorityQueue struct {
	mu sync.Mutex
	cond *sync.Cond
	items itemHeap
	closed bool
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{items: make(itemHeap, 0)}
	pq.cond = sync.NewCond(&pq.mu)
	return pq
}

func (pq *PriorityQueue) Push(j entities.Job) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	if pq.closed {
		return
	}

	heap.Push(&pq.items, &item{job: j})
	pq.cond.Signal()
}

func (pq *PriorityQueue) Pop() (job entities.Job, ok bool) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	for len(pq.items) == 0 && !pq.closed {
		pq.cond.Wait()
	}

	if len(pq.items) == 0 && pq.closed {
		return entities.Job{}, false
	}

	it := heap.Pop(&pq.items).(*item)
	return it.job, true
}

func (pq *PriorityQueue) Len() int {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return len(pq.items)
}

func (pq *PriorityQueue) Close() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	pq.closed = true
	pq.cond.Broadcast()
}

func (pq *PriorityQueue) PromoteExpired(deadlines map[int]time.Duration, maxPriority int) int {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	now := time.Now()
	promoted := 0

	for _, it := range pq.items {
		if it.promoted {
			continue
		}
		deadline, ok := deadlines[it.job.Priority]
		if !ok || deadline <= 0 {
			continue
		}
		if now.Sub(it.job.CreatedAt) >= deadline {
			it.job.Priority = maxPriority
			it.promoted = true
			promoted++
		}
	}

	if promoted > 0 {
		heap.Init(&pq.items)
		pq.cond.Broadcast()
	}

	return promoted
}