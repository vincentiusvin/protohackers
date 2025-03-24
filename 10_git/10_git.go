package main

import (
	"bufio"
	"log"
	"net"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	rw := bufio.NewReadWriter(r, w)
	addr := c.RemoteAddr().String()
	handleIO(rw, addr)
}
