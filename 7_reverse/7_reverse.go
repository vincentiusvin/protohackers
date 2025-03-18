package main

import (
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
		log.Println("connected")
		go handleSession(session)
	}
}

func handleSession(s *lrcp.LRCPSession) {
	incoming, err := s.Resolve()
	if err != nil {
		return
	}

	outch := make(chan string)

	go func() {
		curr := ""
		for data := range incoming {
			for _, char := range data {
				if char == '\n' {
					outch <- curr
					curr = ""
				} else {
					curr += string(char)
				}
			}
		}
		close(outch)
	}()

	for t := range outch {
		log.Println("Received string", t)
		res := ""
		for _, ch := range t {
			res = string(ch) + res
		}
		s.SendData(res + "\n")
	}
}
