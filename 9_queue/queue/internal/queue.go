package internal

import (
	"container/heap"
	"encoding/json"
	"fmt"
)

type JobStatus string

type Job struct {
	Id    int
	Job   json.RawMessage
	Prio  int
	Queue string

	// client id currently working the job.
	// obtain from GetClientID for every session.
	// cannot be zero.
	ClientID int
}

// data structure for efficient querying of jobs

type Queue struct {
	Name string

	jobs      map[int]*Job
	availJobs *priorityQueue
}

func NewQueue(name string) *Queue {
	return &Queue{
		Name:      name,
		jobs:      make(map[int]*Job),
		availJobs: makePriorityQueue(),
	}
}

func (q *Queue) AddJob(nj *Job) {
	q.jobs[nj.Id] = nj
	heap.Push(q.availJobs, nj)
}

func (q *Queue) GetJob(id int) *Job {
	return q.jobs[id]
}

func (q *Queue) GetPrioJob() *Job {
	maxJob := q.availJobs.Peek()
	return maxJob
}

func (q *Queue) AbortJob(jobId int) {
	j := q.GetJob(jobId)
	j.ClientID = 0

	heap.Push(q.availJobs, j)
}

func (q *Queue) PopJob(jobId int, clientId int) {
	j_raw := heap.Pop(q.availJobs)

	j := j_raw.(*Job)
	if j.Id != jobId {
		err := fmt.Errorf("expected pop to yield %v got %v", jobId, j.Id)
		panic(err)
	}

	j.ClientID = clientId
}

func (q *Queue) DeleteJob(id int) {
	for i, c := range q.availJobs.jobs {
		if c.Id == i {
			heap.Remove(q.availJobs, i)
			break
		}
	}

	delete(q.jobs, id)
}

func (q *Queue) DisconnectAllFrom(clientId int) []int {
	dced := make([]int, 0)
	for _, c := range q.jobs {
		if c.ClientID != clientId {
			continue
		}
		q.AbortJob(c.Id)
		dced = append(dced, c.Id)
	}
	return dced
}
