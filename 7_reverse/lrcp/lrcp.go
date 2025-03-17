package lrcp

import (
	"log"
	"net"
)

// provide a tcp like interface
// LRCPListen(c *net.UDPConn) *LRCPServer
// LRCP.Accept() returns a session
func LRCPListen(c *net.UDPConn) *LRCPServer {
	srv := makeLRCPServer()
	go func() {
		for {
			b := make([]byte, 1000)
			n, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			log.Println("recv from: ", addr)

			str := string(b[n])

			parsed, err := Parse(str)
			if err != nil {
				log.Printf("skipped packet from: %v. err: %v", addr, err)
			}

			sess := parsed.GetSession()
			srv.processSession(sess)
		}
	}()
	return srv
}

type LRCPServer struct {
	sessions map[uint]chan LRCPPackets // otherwise existing packets will continue to be listened by this guy
	newch    chan uint                 // if we have new channels, it will be broadcasted to this guy
}

func makeLRCPServer() *LRCPServer {
	serv := &LRCPServer{
		sessions: make(map[uint]chan LRCPPackets),
		newch:    make(chan uint),
	}
	return serv
}

func (ls *LRCPServer) Accept() {
	sid := <-ls.newch
	ls.sessions[sid] = make(chan LRCPPackets)
}

func (ls *LRCPServer) processSession(sid uint) {
	if ls.sessions[sid] != nil {
		return
	}
	ls.newch <- sid
}
