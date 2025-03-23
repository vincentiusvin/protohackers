package queue

import "encoding/json"

type ResponseStatus string

var (
	StatusOK    ResponseStatus = "ok"
	StatusError ResponseStatus = "error"
	StatusNoJob ResponseStatus = "no-job"
)

type RequestType string

var (
	RequestPut    RequestType = "put"
	RequestGet    RequestType = "get"
	RequestDelete RequestType = "delete"
	RequestAbort  RequestType = "abort"
)

type PutRequest struct {
	Request RequestType     `json:"request"`
	Queue   string          `json:"queue"`
	Pri     int             `json:"pri"`
	Job     json.RawMessage `json:"job"`

	respCh chan *PutResponse
}

type PutResponse struct {
	Status ResponseStatus `json:"status"`
	Id     int            `json:"id"`
}

func (pr *PutRequest) init() {
	pr.respCh = make(chan *PutResponse)
}

type GetRequest struct {
	Request RequestType `json:"request"`
	Queues  []string    `json:"queues"`
	Wait    bool        `json:"wait"` // optional, but zero value is fine here because only true is valid

	ClientID int

	respCh chan *GetResponse
}

type GetResponse struct {
	Status ResponseStatus   `json:"status"`
	Id     *int             `json:"id,omitempty"`
	Job    *json.RawMessage `json:"job,omitempty"`
	Pri    *int             `json:"pri,omitempty"`
	Queue  *string          `json:"queue,omitempty"`
}

func (gr *GetRequest) init() {
	gr.respCh = make(chan *GetResponse)
}

type AbortRequest struct {
	Request RequestType `json:"request"`
	Id      int         `json:"id"`

	ClientID int

	respCh chan *AbortResponse
}

type AbortResponse struct {
	Status ResponseStatus `json:"status"`
}

func (ar *AbortRequest) init() {
	ar.respCh = make(chan *AbortResponse)
}

type DeleteRequest struct {
	Request RequestType `json:"request"`
	Id      int         `json:"id"`

	respCh chan *DeleteResponse
}

type DeleteResponse struct {
	Status ResponseStatus `json:"status"`
}

func (ar *DeleteRequest) init() {
	ar.respCh = make(chan *DeleteResponse)
}
