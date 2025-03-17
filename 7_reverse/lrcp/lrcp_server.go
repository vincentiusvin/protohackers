package lrcp

import (
	"log"
	"net"
)

// provides a tcp like interface
// LRCPListen(c *net.UDPConn) *LRCPServer
// LRCP.Accept() returns a session
type LRCPServer struct {
	newSession chan *LRCPSession     // if we have new sessions, it will be broadcasted to this guy
	sessions   map[uint]*LRCPSession // existing sessions will be handled by the values here
}

func MakeLRCPServer() *LRCPServer {
	serv := &LRCPServer{
		sessions:   make(map[uint]*LRCPSession),
		newSession: make(chan *LRCPSession),
	}
	return serv
}

// Accept a new connection
func (ls *LRCPServer) Accept() *LRCPSession {
	inc := <-ls.newSession
	return inc
}

func (ls *LRCPServer) ListenUDP(c *net.UDPConn) {
	go func() {
		for {
			b := make([]byte, 1000)
			n, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			request := string(b[n])
			response := func(b []byte) error {
				_, err := c.WriteToUDP(b, addr)
				return err
			}

			ls.process(request, response)
			log.Println("recv from: ", addr)
		}
	}()
}

func (ls *LRCPServer) process(request string, response func(b []byte) error) {
	parsed, err := Parse(request)
	if err != nil {
		log.Printf("skipped packet %v. err: %v", request, err)
	}
	sid := parsed.GetSession()
	if ls.sessions[sid] != nil {
		ls.sessions[sid].process(parsed)
	} else {
		session := makeLRCPSession(sid, ls, response)
		ls.sessions[sid] = session
		ls.newSession <- session
	}
}
