package main

import (
	"encoding/binary"
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
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	bytes := []byte{}
	applyParser(parseUint8, bytes)
}

type ParseFunc[T any] func(b []byte) (T, int)

// apply parseFunc that satisfies the definition above,
// return a new byte slice that has been advanced
func applyParser[T any](parseFunc ParseFunc[T], value []byte) (T, []byte) {
	out, n := parseFunc(value)
	return out, value[n:]
}

// Consumes tokens from b to produce a string
// Returns number of bytes consumed and the final string
func parseString(b []byte) (string, int) {
	len := int(b[0])
	str := string(b[1 : len+1])

	return str, len + 1
}

// parse uint16
func parseUint8(b []byte) (uint8, int) {
	return uint8(b[0]), 1
}

// parse uint16
func parseUint16(b []byte) (uint16, int) {
	return binary.BigEndian.Uint16(b), 2
}

// parse uint32
func parseUint32(b []byte) (uint32, int) {
	return binary.BigEndian.Uint32(b), 4
}
