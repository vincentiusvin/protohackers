package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	ch := MakeChatroom()

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		log.Println(c.RemoteAddr(), "connected")
		go ch.handleConnection(c)
	}
}

func validateName(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		uppercase := (c >= 'A') && (c <= 'Z')
		lowercase := (c >= 'a') && (c <= 'z')
		digits := (c >= '0') && (c <= '9')

		if uppercase || lowercase || digits {
			continue
		}

		return false
	}

	return true
}

func (ch *Chatroom) handleConnection(c net.Conn) {
	defer c.Close()

	m := &Member{
		conn: c,
	}

	m.Send("Welcome to budgetchat! What shall I call you?\n")
	name := m.Recv()
	if !validateName(name) {
		m.Send("Sorry, your name is invalid. Disconnecting now...\n")
		return
	}
	m.name = name

	ch.AddUser(m)
	defer ch.RemoveUser(m)

	sc := bufio.NewScanner(c)
	for sc.Scan() {
		msg := sc.Text()
		ch.SendMessage(m, msg)
	}
}

type Member struct {
	conn net.Conn
	name string
}

func (m *Member) Send(s string) {
	m.conn.Write([]byte(s))
}

func (m *Member) Recv() string {
	sc := bufio.NewScanner(m.conn)
	for sc.Scan() {
		return sc.Text()
	}
	return ""
}

type Chatroom struct {
	members []*Member
	mu      sync.Mutex
}

func MakeChatroom() *Chatroom {
	return &Chatroom{
		members: make([]*Member, 0),
	}
}

func (ch *Chatroom) broadcast(msg string, exception ...*Member) {
	// turn this into a map if performance is bad
	inException := func(toCheck *Member) bool {
		for _, member := range exception {
			if member == toCheck {
				return true
			}
		}
		return false
	}

	for _, m := range ch.members {
		if inException(m) {
			continue
		}
		m.Send(msg)
	}
}

func (ch *Chatroom) SendMessage(sender *Member, msg string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	msgFormatted := fmt.Sprintf("[%v] %v\n", sender.name, msg)
	ch.broadcast(msgFormatted, sender)
}

func (ch *Chatroom) AddUser(m *Member) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	m.Send(ch.userList())
	enterMsg := fmt.Sprintf("* %v has entered the room\n", m.name)
	ch.broadcast(enterMsg)

	ch.members = append(ch.members, m)
	log.Println(enterMsg)
}

func (ch *Chatroom) RemoveUser(m *Member) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	leaveMsg := fmt.Sprintf("* %v has left the room\n", m.name)
	ch.broadcast(leaveMsg, m)
	ch.members = Remove(ch.members, m)
	log.Println(leaveMsg)
}

func (ch *Chatroom) userList() string {
	members := make([]string, 0)
	for _, c := range ch.members {
		members = append(members, c.name)
	}
	members_string := strings.Join(members, ", ")

	return "* The room contains: " + members_string + "\n"
}

func Remove[T comparable](s []T, elem T) []T {
	for i, c := range s {
		if c != elem {
			continue
		}

		var after []T
		if len(s) > i+1 {
			after = s[i+1:]
		}
		return append(s[:i], after...)
	}
	return s
}
