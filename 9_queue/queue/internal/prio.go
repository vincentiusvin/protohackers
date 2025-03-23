package internal

import (
	"container/heap"
)

type priorityQueue struct {
	jobs []*Job
}

func makePriorityQueue() *priorityQueue {
	p := &priorityQueue{
		jobs: make([]*Job, 0),
	}
	heap.Init(p)
	return p
}

func (pq *priorityQueue) Len() int {
	return len(pq.jobs)
}

func (pq *priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.jobs[i].Prio > pq.jobs[j].Prio
}

func (pq *priorityQueue) Push(x any) {
	cast := x.(*Job)
	pq.jobs = append(pq.jobs, cast)
}

func (pq *priorityQueue) Pop() any {
	l := len(pq.jobs)
	item := pq.jobs[l-1]
	pq.jobs = pq.jobs[:l-1]
	return item
}

func (pq *priorityQueue) Swap(i, j int) {
	pq.jobs[i], pq.jobs[j] = pq.jobs[j], pq.jobs[i]
}

var pq *priorityQueue
var _ heap.Interface = pq

func (pq *priorityQueue) Peek() *Job {
	return pq.jobs[0]
}
