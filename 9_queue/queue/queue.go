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
	lkpJob map[int]*Job // a job id cache to quickly fetch jobs.

	putCh   chan *PutRequest
	getCh   chan *GetRequest
	delCh   chan *DeleteRequest
	abortCh chan *AbortRequest
}

func NewJobCenter(ctx context.Context) *JobCenter {
	jc := &JobCenter{
		Queues: make(map[string]*Queue),
		lkpJob: make(map[int]*Job),

		putCh:   make(chan *PutRequest),
		getCh:   make(chan *GetRequest),
		delCh:   make(chan *DeleteRequest),
		abortCh: make(chan *AbortRequest),
	}
	go jc.process(ctx)
	return jc
}

func (jc *JobCenter) Put(pr *PutRequest) *PutResponse {
	pr.init()
	jc.putCh <- pr
	return <-pr.respCh
}

func (jc *JobCenter) Get(gr *GetRequest) {
	jc.getCh <- gr
}

func (jc *JobCenter) Delete(dr *DeleteRequest) {
	jc.delCh <- dr
}

func (jc *JobCenter) Abort(ar *AbortRequest) {
	jc.abortCh <- ar
}

func (jc *JobCenter) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pr := <-jc.putCh:
			jc.processPut(pr)
		case dr := <-jc.delCh:
			jc.processDelete(dr)
		}
	}
}

// return the id of the job
func (jc *JobCenter) processPut(pr *PutRequest) {
	q := jc.getQueue(pr.Queue)
	nj := &Job{
		Id:    rand.Int(),
		Job:   pr.Job,
		Prio:  pr.Pri,
		Queue: pr.Queue,
	}

	q.Jobs = append(q.Jobs, nj)

	pr.respCh <- &PutResponse{
		Status: StatusOK,
		Id:     nj.Id,
	}
}

func (jc *JobCenter) processDelete(dr *DeleteRequest) {
	job := jc.getJob(dr.Id)
	q := jc.getQueue(job.Queue)

	filtered := make([]*Job, 0)
	for _, c := range q.Jobs {
		if c == job {
			continue
		}
		filtered = append(filtered, c)
	}
	q.Jobs = filtered
}

func (jc *JobCenter) getQueue(name string) *Queue {
	if _, ok := jc.Queues[name]; !ok {
		jc.Queues[name] = newQueue(name)
	}
	return jc.Queues[name]
}

func (jc *JobCenter) getJob(id int) *Job {
	return jc.lkpJob[id]
}
