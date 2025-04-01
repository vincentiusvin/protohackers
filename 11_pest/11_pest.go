package main

import (
	"errors"
	"log"
	"net"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/pest"
	"protohackers/11_pest/types"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	pc := pest.NewControllerTCP()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		s := NewSession(c, pc)
		go s.run()
	}
}

type Session struct {
	c  net.Conn
	pc pest.Controller

	// client->us
	helloChan chan types.Hello
	visitChan chan types.SiteVisit
}

func NewSession(c net.Conn, pc pest.Controller) *Session {
	s := &Session{
		c:         c,
		pc:        pc,
		helloChan: make(chan types.Hello),
		visitChan: make(chan types.SiteVisit),
	}

	return s
}

func (s *Session) run() {
	defer s.c.Close()

	addr := s.c.RemoteAddr().String()
	log.Printf("[%v] connected\n", addr)

	go s.runParser()

	s.runHandshake()
	s.runVisit()
}

func (s *Session) runHandshake() {
	helloB := infra.Encode(types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	})
	_, err := s.c.Write(helloB)
	if err != nil {
		return
	}

	select {
	case helloReply := <-s.helloChan:
		if helloReply.Protocol != "pestcontrol" || helloReply.Version != 1 {
			s.sendError(pest.ErrInvalidHandshake)
			return
		}
	case <-s.visitChan:
		s.sendError(pest.ErrInvalidHandshake)
	}
}

func (s *Session) runVisit() {
	for {
		select {
		case v := <-s.visitChan:
			log.Printf("[%v] added visit: %v\n", s.c.RemoteAddr(), v)
			err := s.pc.AddSiteVisit(v)

			if err != nil {
				if errors.Is(err, pest.ErrInvalidSiteVisit) ||
					errors.Is(err, pest.ErrInvalidHandshake) {
					s.sendError(err)
				}
				return
			}
		case <-s.helloChan:
			s.sendError(pest.ErrInvalidSiteVisit)
		}
	}
}

func (s *Session) runParser() {
	var curr []byte
	for {
		b := make([]byte, 1024)
		n, err := s.c.Read(b)
		curr = append(curr, b[:n]...)
		if err != nil {
			break
		}

		for {
			res := infra.Parse(curr)
			if res.Error != nil {
				if errors.Is(res.Error, infra.ErrNotEnough) {
					break
				}
				s.sendError(res.Error)
				return
			}

			switch v := res.Value.(type) {
			case types.Hello:
				s.helloChan <- v
			case types.SiteVisit:
				s.visitChan <- v
			default:
				s.sendError(pest.ErrInvalidData)
				return
			}

			curr = res.Next
		}
	}
}

func (s *Session) sendError(err error) {
	log.Println("sent error")
	errorB := infra.Encode(types.Error{
		Message: err.Error(),
	})
	s.c.Write(errorB)
}
