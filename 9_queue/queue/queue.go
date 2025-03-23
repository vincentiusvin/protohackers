package queue

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
)

type JobStatus string

type Job struct {
	Id    int
	Job   json.RawMessage
	Prio  int
	Queue string

	ClientID int // client id currently working the job
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

	putCh   chan *PutRequest
	getCh   chan *GetRequest
	delCh   chan *DeleteRequest
	abortCh chan *AbortRequest
}

func NewJobCenter(ctx context.Context) *JobCenter {
	jc := &JobCenter{
		Queues: make(map[string]*Queue),

		putCh:   make(chan *PutRequest),
		getCh:   make(chan *GetRequest),
		delCh:   make(chan *DeleteRequest),
		abortCh: make(chan *AbortRequest),
	}
	log.Printf("initialized new job center\n")
	go jc.process(ctx)
	return jc
}

func (jc *JobCenter) Put(pr *PutRequest) *PutResponse {
	pr.init()
	jc.putCh <- pr
	return <-pr.respCh
}

func (jc *JobCenter) Get(gr *GetRequest) *GetResponse {
	gr.init()
	jc.getCh <- gr
	return <-gr.respCh
}

func (jc *JobCenter) Abort(ar *AbortRequest) *AbortResponse {
	ar.init()
	jc.abortCh <- ar
	return <-ar.respCh
}

func (jc *JobCenter) Delete(dr *DeleteRequest) *DeleteResponse {
	dr.init()
	jc.delCh <- dr
	return <-dr.respCh
}

func (jc *JobCenter) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case pr := <-jc.putCh:
			jc.processPut(pr)
		case gr := <-jc.getCh:
			jc.processGet(gr)
		case ar := <-jc.abortCh:
			jc.processAbort(ar)
		case dr := <-jc.delCh:
			jc.processDelete(dr)
		}
	}
}

// return the id of the job
func (jc *JobCenter) processPut(pr *PutRequest) {
	log.Printf("put q:%v pri:%v\n", pr.Queue, pr.Pri)

	var resp PutResponse
	q := jc.getQueue(pr.Queue)
	nj := &Job{
		Id:    rand.Int(),
		Job:   pr.Job,
		Prio:  pr.Pri,
		Queue: pr.Queue,
	}

	resp.Status = StatusOK
	resp.Id = nj.Id

	q.Jobs = append(q.Jobs, nj)

	log.Printf("put ret:%v\n", resp.Status)
	pr.respCh <- &resp
}

func (jc *JobCenter) processGet(gr *GetRequest) {
	log.Printf("get qs:%v w:%v", gr.Queues, gr.Wait)
	var maxJob *Job
	var resp GetResponse

	for _, q := range jc.Queues {
		for _, j := range q.Jobs {
			if maxJob == nil || maxJob.Prio < j.Prio {
				maxJob = j
			}
		}
	}

	if maxJob == nil {
		if !gr.Wait {
			resp.Status = StatusNoJob
		} else {
			panic("wait not implemented yet")
		}
	} else {
		resp.Status = StatusOK
		resp.Id = &maxJob.Id
		resp.Job = &maxJob.Job
		resp.Pri = &maxJob.Prio
		resp.Queue = &maxJob.Queue
	}

	log.Printf("get ret:%v\n", resp.Status)
	gr.respCh <- &resp
}

func (jc *JobCenter) processAbort(ar *AbortRequest) {
	log.Printf("abort id:%v", ar.Id)
	job, _ := jc.findJob(ar.Id)

	var resp AbortResponse

	if job == nil {
		resp.Status = StatusNoJob
	} else if job.ClientID != ar.ClientID {
		resp.Status = StatusError
	} else {
		resp.Status = StatusOK
	}

	log.Printf("abort ret:%v\n", resp.Status)
	ar.respCh <- &resp
}

func (jc *JobCenter) processDelete(dr *DeleteRequest) {
	log.Printf("delete id:%v", dr.Id)
	var resp DeleteResponse

	job, queue := jc.findJob(dr.Id)
	filtered := make([]*Job, 0)
	found := false
	for _, c := range queue.Jobs {
		if c == job {
			found = true
			continue
		}
		filtered = append(filtered, c)
	}
	queue.Jobs = filtered

	if found {
		resp.Status = StatusOK
	} else {
		resp.Status = StatusNoJob
	}

	log.Printf("delete ret:%v", resp.Status)
	dr.respCh <- &resp
}

func (jc *JobCenter) getQueue(name string) *Queue {
	if _, ok := jc.Queues[name]; !ok {
		jc.Queues[name] = newQueue(name)
	}
	return jc.Queues[name]
}

func (jc *JobCenter) findJob(id int) (*Job, *Queue) {
	for _, q := range jc.Queues {
		for _, j := range q.Jobs {
			if j.Id == id {
				return j, q
			}
		}
	}
	return nil, nil
}
