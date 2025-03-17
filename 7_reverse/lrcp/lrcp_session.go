package lrcp

import "fmt"

type LRCPSession struct {
	// user facing
	sid      uint
	resolved chan string

	// server facing
	srv      *LRCPServer
	inbuf    []string
	incoming chan LRCPPackets
	outbuf   []string
	outgoing chan LRCPPackets
}

func makeLRCPSession(
	sid uint,
	incoming chan LRCPPackets,
	outgoing chan LRCPPackets,
	srv *LRCPServer,
) *LRCPSession {
	ret := &LRCPSession{
		sid:      sid,
		resolved: make(chan string),

		srv:      srv,
		inbuf:    make([]string, 0),
		incoming: incoming,
		outbuf:   make([]string, 0),
		outgoing: outgoing,
	}

	go ret.processIncoming()
	go ret.processOutgoing()

	return ret
}

func (ls *LRCPSession) processIncoming() {
	for str := range ls.incoming {
		switch v := str.(type) {
		case *Data:
			fmt.Println(v.Data)
		}
	}
}

func (ls *LRCPSession) processOutgoing() {
}

func (ls *LRCPSession) Send(s string) {
}

func (ls *LRCPSession) Resolve() chan string {
	return ls.resolved
}
