package lrcp

import (
	"log"
	"net"
)

// provide a tcp like interface
// LRCPListen(c *net.UDPConn) *LRCPServer
// LRCP.Accept() returns a session
func LRCPListenUDP(c *net.UDPConn) *LRCPServer {
	srv := MakeLRCPServer()
	ch := make(chan string)
	go func() {
		for {
			b := make([]byte, 1000)
			n, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			log.Println("recv from: ", addr)
			ch <- string(b[n])
		}
	}()
	go srv.Process(ch)
	return srv
}

type LRCPSession struct {
	Sid uint
	inc chan LRCPPackets
	out chan LRCPPackets
	srv *LRCPServer
}

func makeLRCPSession(sid uint, srv *LRCPServer) *LRCPSession {
	return &LRCPSession{
		inc: make(chan LRCPPackets),
		out: make(chan LRCPPackets),
		Sid: sid,
		srv: srv,
	}
}

type LRCPServer struct {
	newch    chan uint             // if we have new channels, it will be broadcasted to this guy
	sessions map[uint]*LRCPSession // existing channels will be handled by the values here
}

func MakeLRCPServer() *LRCPServer {
	serv := &LRCPServer{
		sessions: make(map[uint]*LRCPSession),
		newch:    make(chan uint),
	}
	return serv
}

func (ls *LRCPServer) Accept() *LRCPSession {
	sid := <-ls.newch
	ls.sessions[sid] = makeLRCPSession(sid, ls)
	return ls.sessions[sid]
}

func (ls *LRCPServer) Process(ch chan string) {
	for str := range ch {
		parsed, err := Parse(str)
		if err != nil {
			log.Printf("skipped packet %v. err: %v", str, err)
		}
		sid := parsed.GetSession()
		if ls.sessions[sid] != nil {
			continue
		}
		ls.newch <- sid
	}
}
