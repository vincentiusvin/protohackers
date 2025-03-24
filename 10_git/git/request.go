package git

type Status string

var (
	StatusOK  Status = "OK"
	StatusERR Status = "ERR"
)

type HelpRequest struct{}

type GetRequest struct {
	Filename string // begin with /
	Revision int
}

type PutRequest struct {
	Filename string
	Data     []byte
}

type PutResponse struct {
	Status   Status
	ErrMsg   string
	Filename string
	Revision int
}

type ListRequest struct {
	Dir string
}

// entries are in the format of:
// dirname DIR
// filename r1
type ListResponse struct {
	Status  Status
	ErrMsg  string
	Entries []string
}
