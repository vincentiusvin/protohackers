package main

import (
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
		_, err := ln.Accept()
		if err != nil {
			panic(err)
		}
	}
}
