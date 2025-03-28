package queue

import (
	"context"
	"log"
	"protohackers/9_queue/queue/internal"
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

	q.AddJob(nj)

	log.Printf("put ret:%v\n", resp.Status)
	pr.respCh <- &resp
}

func (jc *JobCenter) processGet(gr *GetRequest) {
	log.Printf("get qs:%v w:%v", gr.Queues, gr.Wait)

	var maxJob *internal.Job
	var maxQueue *internal.Queue
	var resp GetResponse

	queues := make(map[string]bool)
	for _, q := range gr.Queues {
		queues[q] = true
	}

	for _, q := range jc.Queues {
		if !queues[q.Name] {
			continue
		}
		j := q.GetPrioJob()
		if j == nil {
			continue
		}
		if maxJob == nil || maxJob.Prio < j.Prio {
			maxJob = j
			maxQueue = q
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
		maxQueue.PopJob(maxJob.Id, gr.ClientID)

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

	var resp AbortResponse

	defer func() {
		ar.respCh <- &resp
	}()

	q := jc.findQueueFromJob(ar.Id)

	if q == nil {
		resp.Status = StatusNoJob
		return
	}

	job := q.GetJob(ar.Id)

	if job == nil {
		resp.Status = StatusNoJob
		return
	}

	if job.ClientID != ar.ClientID {
		resp.Status = StatusError
		return
	}

	q.AbortJob(ar.Id)
	resp.Status = StatusOK

	log.Printf("abort ret:%v\n", resp.Status)
}

func (jc *JobCenter) processDelete(dr *DeleteRequest) {
	log.Printf("delete id:%v", dr.Id)
	var resp DeleteResponse

	q := jc.findQueueFromJob(dr.Id)

	if q == nil {
		resp.Status = StatusNoJob
	} else {
		q.DeleteJob(dr.Id)
		resp.Status = StatusOK
	}

	log.Printf("delete ret:%v", resp.Status)
	dr.respCh <- &resp
}

func (jc *JobCenter) processDisconnection(dr *DisconnectRequest) {
	aborted := make([]int, 0)
	for _, q := range jc.Queues {
		ab := q.DisconnectAllFrom(dr.ClientID)
		aborted = append(aborted, ab...)
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

func (jc *JobCenter) findQueueFromJob(jobId int) *internal.Queue {
	for _, q := range jc.Queues {
		if q.GetJob(jobId) != nil {
			return q
		}
	}
	return nil
}

// get a monotonically increasing job id.
// will never return zero
func (jc *JobCenter) getJobID() int {
	jc.currJobId += 1
	return jc.currJobId
}
