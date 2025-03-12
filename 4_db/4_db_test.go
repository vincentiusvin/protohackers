package main

import (
	"fmt"
	"net"
	"testing"
)

func TestDial(t *testing.T) {
	c, err := net.ListenPacket("udp", ":8000")
	if err != nil {
		panic(err)
	}
	for {
		b := make([]byte, 1000)
		_, addr, err := c.ReadFrom(b)
		if err != nil {
			panic(err)
		}
		fmt.Println(addr.Network(), addr.String())

		// This sends the reply on another port
		// (debugged using tcpdump)
		nc, err := net.Dial(addr.Network(), addr.String())
		if err != nil {
			panic(err)
		}

		n, err := nc.Write([]byte("oy"))
		if err != nil {
			panic(err)
		}
		fmt.Println(n)
		nc.Close()
	}
}

func TestNoDial(t *testing.T) {
	c, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 8000,
	})

	if err != nil {
		panic(err)
	}
	defer c.Close()

	for {
		b := make([]byte, 1000)
		_, addr, err := c.ReadFromUDP(b)
		if err != nil {
			panic(err)
		}

		// This sends the reply on the same server port (8000)
		// (debugged using tcpdump)
		n, err := c.WriteToUDP([]byte("oy"), addr)
		if err != nil {
			panic(err)
		}
		fmt.Println(n)
	}
}
