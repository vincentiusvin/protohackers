package main

import (
	"bufio"
	"log"
	"net"
	"protohackers/10_git/git"
	"protohackers/10_git/handler"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	vc := git.NewVersionControl()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c, vc)
	}
}

func handleConn(c net.Conn, vc *git.VersionControl) {
	defer c.Close()

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	rw := bufio.NewReadWriter(r, w)
	addr := c.RemoteAddr().String()
	handler.HandleIO(rw, addr, vc)
}
