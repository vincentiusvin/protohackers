package main

import (
	"fmt"
	"io"
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
		go Echo3(c)
	}
}

// naive implementation with for loops
func Echo1(c net.Conn) {
	defer c.Close()
	for {
		sl := make([]byte, 64)
		n, err := c.Read(sl)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Got %v bytes\n", n)
		n, err = c.Write(sl)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Printf("Sent %v bytes\n", n)
	}
}

// chans
func Echo2(c net.Conn) {
	cha := make(chan []byte)
	defer c.Close()

	go func() {
		for {
			b := make([]byte, 64)
			_, err := c.Read(b)
			if err != nil {
				close(cha)
				break
			}
			cha <- b
		}
	}()

	for b := range cha {
		_, err := c.Write(b)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// lol
func Echo3(c net.Conn) {
	io.Copy(c, c)
	defer c.Close()
}
