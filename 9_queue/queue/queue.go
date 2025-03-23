package queue

import (
	"context"
	"math/rand"
)

type Job struct {
	Id    int
	Job   any
	Prio  int
	Queue string
}

type Queue struct {
	Name string
	Jobs []*Job // sorted descendingly by prio. use heap if it needs to be faster
}

func newQueue(name string) *Queue {
	return &Queue{
		Name: name,
		Jobs: make([]*Job, 0),
	}
}

type JobCenter struct {
	Queues map[string]*Queue
	putCh  chan PutRequest
}

func NewJobCenter(ctx context.Context) *JobCenter {
	jc := &JobCenter{
		putCh: make(chan PutRequest),
	}
	go jc.Process(ctx)
	return jc
}

func (jc *JobCenter) Process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pr := <-jc.putCh:
			jc.processPut(pr)
		}
	}
}

// return the id of the job
func (jc *JobCenter) processPut(pr PutRequest) int {
	q := jc.getQueue(pr.Queue)
	nj := &Job{
		Id:    rand.Int(),
		Job:   pr.Job,
		Prio:  pr.Pri,
		Queue: pr.Queue,
	}

	q.Jobs = append(q.Jobs, nj)

	return nj.Id
}

func (jc *JobCenter) getQueue(name string) *Queue {
	if _, ok := jc.Queues[name]; !ok {
		jc.Queues[name] = newQueue(name)
	}
	return jc.Queues[name]
}
