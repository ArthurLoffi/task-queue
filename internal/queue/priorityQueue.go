package queue

import (
	"container/heap"
	"sync"
	"task-queue/internal/entities"
)

type item struct {
	job entities.Job
	index int
}

type itemHeap []*item

func (h itemHeap) Len() int {
	return len(h)
}

func (h itemHeap) Less(i, j int) bool {
	return h[i].job.Priority > h[j].job.Priority
}

func (h itemHeap) Swap(i, j int) {
	h[i], h[j] = h[i], h[j]
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

func (pq *PriorityQueue) Pop() {}

func (pq *PriorityQueue) Len() {}

func (pq *PriorityQueue) Close() {}