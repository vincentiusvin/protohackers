package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		echo3(c)
	}
}

// naive implementation with for loops
func echo1(c net.Conn) {
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
func echo2(c net.Conn) {
	cha := make(chan []byte)

	go func() {
		for {
			b := make([]byte, 64)
			_, err := c.Read(b)
			if err != nil {
				close(cha)
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

func echo3(c net.Conn) {
	io.Copy(c, c)
}
