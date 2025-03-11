package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		prime(c)
	}
}

func prime(c net.Conn) {
	defer c.Close()
	sc := bufio.NewScanner(c)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
}
