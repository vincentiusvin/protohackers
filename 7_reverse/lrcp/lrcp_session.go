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
	case *Connect:
		ls.handleConnect()
	case *Data:
		ls.handleData(v)
	case *Close:
		ls.handleClose()
	}
}

func (ls *LRCPSession) handleConnect() {
	ack := &Ack{
		Session: ls.sid,
		Length:  uint(len(ls.buff)),
	}
	ls.Send(ack.Encode())
}

func (ls *LRCPSession) handleData(v *Data) {
	if v.Pos == uint(len(ls.buff)) {
		ls.buff += v.Data
		ls.resolved <- v.Data
	}
	ack := &Ack{
		Session: ls.sid,
		Length:  uint(len(ls.buff)),
	}
	ls.Send(ack.Encode())
}

func (ls *LRCPSession) handleClose() {
	close(ls.resolved)
	ack := &Close{
		Session: ls.sid,
	}
	ls.Send(ack.Encode())
}

func (ls *LRCPSession) Send(s string) {
	ls.response([]byte(s))
}

func (ls *LRCPSession) Resolve() chan string {
	return ls.resolved
}
