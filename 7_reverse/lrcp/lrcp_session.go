package lrcp

type LRCPSession struct {
	sid      uint                 // session id
	resolved chan string          // output the ordered string
	srv      *LRCPServer          // the server that handles us
	response func(b []byte) error // callback to reply to messages
	buff     string
}

func makeLRCPSession(
	sid uint,
	srv *LRCPServer,
	response func(b []byte) error,
) *LRCPSession {
	return &LRCPSession{
		sid:      sid,
		resolved: make(chan string),
		srv:      srv,
		response: response,
	}
}

func (ls *LRCPSession) process(p LRCPPackets) {
	switch v := p.(type) {
	case *Data:
		if v.Pos == uint(len(ls.buff)) {
			ls.buff += v.Data
			ls.resolved <- v.Data
		}
	case *Close:
		ls.Close(v)
	}
}

func (ls *LRCPSession) Close(v *Close) {
	close(ls.resolved)
	ls.Send(v.Encode())
}

func (ls *LRCPSession) Send(s string) {
	ls.response([]byte(s))
}

func (ls *LRCPSession) Resolve() chan string {
	return ls.resolved
}
