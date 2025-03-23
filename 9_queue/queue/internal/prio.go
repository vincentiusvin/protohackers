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

type Queue struct {
	Name string
	Jobs []*Job
}

func NewQueue(name string) *Queue {
	return &Queue{
		Name: name,
		Jobs: make([]*Job, 0),
	}
}
