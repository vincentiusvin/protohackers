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
}

type PutResponse struct {
	Status ResponseStatus `json:"status"`
	Id     int            `json:"id"`
}

type GetRequest struct {
	Request RequestType `json:"request"`
	Queues  []string    `json:"queues"`
	Wait    bool        `json:"wait"` // optional, but zero value is fine here because only true is valid
}

type GetResponse struct {
	Status ResponseStatus   `json:"status"`
	Id     *int             `json:"id,omitempty"`
	Queue  *string          `json:"queue,omitempty"`
	Pri    *int             `json:"pri,omitempty"`
	Job    *json.RawMessage `json:"job,omitempty"`
}

type AbortRequest struct {
	Request RequestType `json:"request"`
	Id      int         `json:"id"`
}

type AbortResponse struct {
	Status ResponseStatus `json:"status"`
}

type DeleteRequest struct {
	Request RequestType `json:"request"`
	Id      int         `json:"id"`
}

type DeleteResponse struct {
	Status ResponseStatus `json:"status"`
}
