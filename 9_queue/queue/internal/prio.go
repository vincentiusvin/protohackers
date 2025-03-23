package internal

import "encoding/json"

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

	jobs map[int]*Job
}

func NewQueue(name string) *Queue {
	return &Queue{
		Name: name,
		jobs: make(map[int]*Job),
	}
}

func (q *Queue) AddJob(nj *Job) {
	q.jobs[nj.Id] = nj
}

func (q *Queue) GetJob(id int) *Job {
	for _, c := range q.jobs {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (q *Queue) GetPrioJob() *Job {
	var maxJob *Job

	for _, c := range q.jobs {
		if c.ClientID != 0 {
			continue
		}
		if maxJob == nil || c.Prio > maxJob.Prio {
			maxJob = c
		}
	}

	return maxJob
}

func (q *Queue) AbortJob(jobId int) {
	j := q.GetJob(jobId)
	j.ClientID = 0
}

func (q *Queue) MarkJobExecuting(jobId int, clientId int) {
	j := q.GetJob(jobId)
	j.ClientID = clientId
}

func (q *Queue) DeleteJob(id int) {
	delete(q.jobs, id)
}

func (q *Queue) DisconnectAllFrom(clientId int) []int {
	dced := make([]int, 0)
	for _, c := range q.jobs {
		if c.ClientID != clientId {
			continue
		}
		c.ClientID = 0
		dced = append(dced, c.Id)
	}
	return dced
}
