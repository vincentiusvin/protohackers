package lrcp

import (
	"fmt"
	"log"
	"time"
)

var (
	errClosedSession = fmt.Errorf("lrcp session closed")
)

type LRCPSession struct {
	closed   bool
	sid      uint                 // session id
	resolved chan string          // output the ordered string
	srv      *LRCPServer          // the server that handles us
	response func(b []byte) error // callback to reply to messages
	buff     string

	sendCh chan string
	ackCh  chan uint
}

func makeLRCPSession(
	sid uint,
	srv *LRCPServer,
	response func(b []byte) error,
) *LRCPSession {
	s := &LRCPSession{
		sid:      sid,
		resolved: make(chan string),
		srv:      srv,
		response: response,
		sendCh:   make(chan string),
		ackCh:    make(chan uint),
	}
	go s.runSender()
	return s
}

// send data over the LRCP connection
func (ls *LRCPSession) SendData(s string) error {
	if ls.closed {
		return errClosedSession
	}

	ls.sendCh <- s

	return nil
}

// get the channel that outputs ordered packets
func (ls *LRCPSession) Resolve() (chan string, error) {
	if ls.closed {
		return nil, errClosedSession
	}
	return ls.resolved, nil
}

func (ls *LRCPSession) sendRaw(s string) error {
	return ls.response([]byte(s))
}

func (ls *LRCPSession) handlePacket(p LRCPPackets) {
	switch v := p.(type) {
	case *Connect:
		ls.handleConnect()
	case *Data:
		ls.handleData(v)
	case *Close:
		ls.handleClose()
	case *Ack:
		ls.ackCh <- v.Length
	}
}

func (ls *LRCPSession) handleConnect() {
	ack := &Ack{
		Session: ls.sid,
		Length:  uint(len(ls.buff)),
	}
	ls.sendRaw(ack.Encode())
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
	ls.sendRaw(ack.Encode())
}

func (ls *LRCPSession) handleClose() {
	if ls.closed {
		return
	}
	close(ls.resolved)
	ack := &Close{
		Session: ls.sid,
	}
	ls.sendRaw(ack.Encode())
	ls.srv.Unregister(ls)
	ls.closed = true
}

func (ls *LRCPSession) runSender() {
	toSend := ""
	var sent uint = 0
	startSending := make(chan struct{}, 1)

	for !ls.closed {
		select {
		case send := <-ls.sendCh:
			toSend += send
			// probably can code this using a condvar
			select {
			case startSending <- struct{}{}:
			default:
			}
		case ack := <-ls.ackCh:
			maxLen := uint(len(toSend))
			if ack > maxLen {
				ls.handleClose()
				return
			}
			if ack > sent {
				sent = ack
				log.Printf("advancing sent to %v out of %v\n", sent, maxLen)
			}
			if sent != maxLen { // immediately send the next one
				select {
				case startSending <- struct{}{}:
				default:
				}
			}
		case <-startSending:
			currLen := uint(len(toSend))
			sendingLen := min(sent+800, currLen)
			if sent == sendingLen {
				continue
			}
			forward := toSend[sent:sendingLen]
			data := &Data{
				Session: ls.sid,
				Pos:     sent,
				Data:    forward,
			}
			err := ls.sendRaw(data.Encode())
			log.Printf("sending bytes %v..%v out of %v\n", sent, sendingLen, currLen)
			if err != nil {
				log.Printf("failed to send: %v\n", err)
			}
		// the website told us to use 3 seconds but it's too long.
		// the 60 second timeout will kick in if the packet loss is too high.
		case <-time.After(1 * time.Second):
			select {
			case startSending <- struct{}{}:
			default:
			}
		}
	}
}
