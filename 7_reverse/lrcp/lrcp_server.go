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
	ls.Listen(func() LRCPListenerSession {
		return &LRCPUDPSession{c: c}
	})
}

func (ls *LRCPServer) Unregister(s *LRCPSession) {
	delete(ls.sessions, s.sid)
}

type LRCPListenerSession interface {
	Read() ([]byte, error)
	Write([]byte) error
}

type LRCPListener func() LRCPListenerSession

func (ls *LRCPServer) Listen(list LRCPListener) {
	go func() {
		for {
			sess := list()
			b, err := sess.Read()
			if err != nil {
				panic(err)
			}
			request := string(b)
			ls.process(request, sess.Write)
		}
	}()
}

func (ls *LRCPServer) process(request string, response func(b []byte) error) {
	parsed, err := Parse(request)
	if err != nil {
		log.Printf("skipped packet %v. err: %v", request, err)
		return
	}
	sid := parsed.GetSession()
	if ls.sessions[sid] != nil {
		ls.sessions[sid].handlePacket(parsed)
	} else {
		session := makeLRCPSession(sid, ls, response)
		ls.newSession <- session // make sure it is accepted first before registering

		ls.sessions[sid] = session
		session.handlePacket(parsed)
	}
}

type LRCPUDPSession struct {
	c    *net.UDPConn
	addr *net.UDPAddr
}

func (lus *LRCPUDPSession) Read() ([]byte, error) {
	b := make([]byte, 1000)
	n, addr, err := lus.c.ReadFromUDP(b)
	if err != nil {
		return nil, err
	}
	lus.addr = addr
	return b[:n], nil
}

func (lus *LRCPUDPSession) Write(b []byte) error {
	if lus.addr == nil {
		panic("need to call read first before write")
	}

	_, err := lus.c.WriteToUDP(b, lus.addr)
	return err
}
