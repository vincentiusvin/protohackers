package lrcp

import "fmt"

type LRCPSession struct {
	sid      uint                 // session id
	resolved chan string          // output the ordered string
	srv      *LRCPServer          // the server that handles us
	response func(b []byte) error // callback to reply to messages
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
		fmt.Println(v.Data)
	}
}

func (ls *LRCPSession) Send(s string) {
}

func (ls *LRCPSession) Resolve() chan string {
	return ls.resolved
}
