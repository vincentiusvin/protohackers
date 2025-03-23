package queue

import (
	"context"
	"fmt"
	"log"
	"protohackers/9_queue/queue/internal"
	"time"
)

type JobCenter struct {
	Queues map[string]*internal.Queue

	putCh   chan *PutRequest
	getCh   chan *GetRequest
	delCh   chan *DeleteRequest
	abortCh chan *AbortRequest
	discCh  chan *DisconnectRequest

	currJobId    int // monotonic job id
	currClientID int // monotonic client id

	waitingRequests []*GetRequest
}

func NewJobCenter(ctx context.Context) *JobCenter {
	jc := &JobCenter{
		Queues: make(map[string]*internal.Queue),

		putCh:   make(chan *PutRequest),
		getCh:   make(chan *GetRequest),
		delCh:   make(chan *DeleteRequest),
		abortCh: make(chan *AbortRequest),
		discCh:  make(chan *DisconnectRequest),
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

func (jc *JobCenter) DisconnectWorker(dr *DisconnectRequest) {
	dr.init()
	jc.discCh <- dr
	<-dr.respCh
}

// get a monotonically increasing client id.
// will never return zero
func (jc *JobCenter) GetClientID() int {
	jc.currClientID += 1
	return jc.currClientID
}

func (jc *JobCenter) process(ctx context.Context) {
	for {
		processWaiting := true
		t := time.Now()
		select {
		case <-ctx.Done():
			return
		case pr := <-jc.putCh:
			jc.processPut(pr)
		case gr := <-jc.getCh:
			jc.processGet(gr)
			processWaiting = false // dont need to handle waiting requests if last was a get
		case ar := <-jc.abortCh:
			jc.processAbort(ar)
		case dr := <-jc.delCh:
			jc.processDelete(dr)
		case dc := <-jc.discCh:
			jc.processDisconnection(dc)
		}
		log.Printf("dur: %v\n", time.Since(t))
		if processWaiting {
			jc.processWaitingRequests()
		}
	}
}

// return the id of the job
func (jc *JobCenter) processPut(pr *PutRequest) {
	log.Printf("put q:%v pri:%v\n", pr.Queue, pr.Pri)

	var resp PutResponse
	q := jc.getQueue(pr.Queue)
	nj := &internal.Job{
		Id:    jc.getJobID(),
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
	var maxJob *internal.Job
	var resp GetResponse

	queues := make(map[string]bool)
	for _, q := range gr.Queues {
		queues[q] = true
	}

	for _, q := range jc.Queues {
		if !queues[q.Name] {
			continue
		}
		for _, j := range q.Jobs {
			if j.ClientID != 0 {
				continue
			}
			if maxJob == nil || maxJob.Prio < j.Prio {
				maxJob = j
			}
		}
	}

	if maxJob == nil {
		if gr.Wait {
			jc.waitingRequests = append(jc.waitingRequests, gr)
			return
		}
		resp.Status = StatusNoJob
		log.Printf("get ret:%v\n", resp.Status)
	} else {
		maxJob.ClientID = gr.ClientID

		resp.Status = StatusOK
		resp.Id = &maxJob.Id
		resp.Job = &maxJob.Job
		resp.Pri = &maxJob.Prio
		resp.Queue = &maxJob.Queue
		log.Printf("get ret:%v id:%v prio:%v q:%v\n", resp.Status, maxJob.Id, maxJob.Prio, maxJob.Queue)
	}

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
		job.ClientID = 0
		resp.Status = StatusOK
	}

	log.Printf("abort ret:%v\n", resp.Status)
	ar.respCh <- &resp
}

func (jc *JobCenter) processDelete(dr *DeleteRequest) {
	log.Printf("delete id:%v", dr.Id)
	var resp DeleteResponse

	job, queue := jc.findJob(dr.Id)

	if job != nil && queue != nil {
		filtered := make([]*internal.Job, 0)
		found := false
		for _, c := range queue.Jobs {
			if c == job {
				found = true
				continue
			}
			filtered = append(filtered, c)
		}
		queue.Jobs = filtered
		if !found {
			panic(fmt.Sprintf("expected job %v to be in queue %v", job, queue))
		}
		resp.Status = StatusOK
	} else {
		resp.Status = StatusNoJob

	}

	log.Printf("delete ret:%v", resp.Status)
	dr.respCh <- &resp
}

func (jc *JobCenter) processDisconnection(dr *DisconnectRequest) {
	aborted := make([]int, 0)
	for _, q := range jc.Queues {
		for _, j := range q.Jobs {
			if j.ClientID == dr.ClientID {
				j.ClientID = 0
				aborted = append(aborted, j.Id)
			}
		}
	}
	log.Printf("disconnected %v. aborted %v", dr.ClientID, aborted)
	dr.respCh <- struct{}{}
}

func (jc *JobCenter) processWaitingRequests() {
	if jc.waitingRequests == nil {
		return
	}

	processing := jc.waitingRequests
	jc.waitingRequests = nil

	for _, c := range processing {
		jc.processGet(c) // this will readd to waitingRequests if still unresolved
	}
}

func (jc *JobCenter) getQueue(name string) *internal.Queue {
	if _, ok := jc.Queues[name]; !ok {
		jc.Queues[name] = internal.NewQueue(name)
	}
	return jc.Queues[name]
}

func (jc *JobCenter) findJob(id int) (*internal.Job, *internal.Queue) {
	for _, q := range jc.Queues {
		for _, j := range q.Jobs {
			if j.Id == id {
				return j, q
			}
		}
	}
	return nil, nil
}

// get a monotonically increasing job id.
// will never return zero
func (jc *JobCenter) getJobID() int {
	jc.currJobId += 1
	return jc.currJobId
}
