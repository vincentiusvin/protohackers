package lrcp

import (
	"fmt"
	"log"
	"net"
)

type lrcpIncoming struct {
	msg   string
	reply func(string)
}

// provide a tcp like interface
// LRCPListen(c *net.UDPConn) *LRCPServer
// LRCP.Accept() returns a session
func LRCPListenUDP(c *net.UDPConn) *LRCPServer {
	srv := MakeLRCPServer()
	ch := make(chan lrcpIncoming)
	go func() {
		for {
			b := make([]byte, 1000)
			n, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			ch <- lrcpIncoming{
				msg: string(b[n]),
				reply: func(s string) {
					c.WriteToUDP([]byte(s), addr)
				},
			}

			log.Println("recv from: ", addr)
		}
	}()
	go srv.Process(ch)
	return srv
}

type LRCPServer struct {
	newch    chan lrcpIncoming     // if we have new sessions, it will be broadcasted to this guy
	sessions map[uint]*LRCPSession // otherwise, existing sessions will be handled by the values here
}

func MakeLRCPServer() *LRCPServer {
	serv := &LRCPServer{
		sessions: make(map[uint]*LRCPSession),
		newch:    make(chan lrcpIncoming),
	}
	return serv
}

// Accept a new connection
func (ls *LRCPServer) Accept() (*LRCPSession, error) {
	inc := <-ls.newch
	psed, err := Parse(inc.msg)
	if err != nil {
		return nil, fmt.Errorf("cannot accept: %w", err)
	}
	conn := psed.(*Connect)
	sid := conn.GetSession()

	servToSess := make(chan LRCPPackets)
	sessToServ := make(chan LRCPPackets)

	ls.sessions[sid] = makeLRCPSession(sid, servToSess, sessToServ, ls)
	return ls.sessions[sid], nil
}

// Execute to start listening for requests
func (ls *LRCPServer) Process(ch chan lrcpIncoming) {
	for inc := range ch {
		str := inc.msg
		parsed, err := Parse(str)
		if err != nil {
			log.Printf("skipped packet %v. err: %v", str, err)
		}
		sid := parsed.GetSession()
		if ls.sessions[sid] != nil {
			ls.sessions[sid].incoming <- parsed
		} else {
			ls.newch <- inc
		}
	}
}
