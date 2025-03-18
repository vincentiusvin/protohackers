package main

import (
	"fmt"
	"log"
	"net"
	"protohackers/7_reverse/lrcp"
	"strconv"
)

func main() {
	port := 8000
	c, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at port", strconv.Itoa(port))

	defer c.Close()

	serv := lrcp.MakeLRCPServer()
	serv.ListenUDP(c)

	for {
		session := serv.Accept()
		handleSession(session)
	}
}

func handleSession(s *lrcp.LRCPSession) {
	incoming, err := s.Resolve()
	if err != nil {
		return
	}

	for data := range incoming {
		fmt.Println(data)
	}
}
