package main

import (
	"encoding/binary"
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
		go means(c)
	}
}

type Insert struct {
	Timestamp int
	Price     int
}

type Query struct {
	Mintime int
	Maxtime int
}

func parsePacket(b []byte) (*Insert, *Query) {
	t := b[0]

	// uint32 -> int32 -> int
	// casting directly to int will break negative numbers on 64 bit systems
	firstint := b[1:5]
	firstnum := int(int32(binary.BigEndian.Uint32(firstint)))

	secint := b[5:9]
	secnum := int(int32(binary.BigEndian.Uint32(secint)))

	if t == 'Q' {
		return nil, &Query{
			Mintime: firstnum,
			Maxtime: secnum,
		}
	} else if t == 'I' {
		return &Insert{
			Timestamp: firstnum,
			Price:     secnum,
		}, nil
	}
	return nil, nil
}

func means(c net.Conn) {
	b := make([]byte, 9)
	c.Read(b)
}
