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
}

type Plate struct {
	Plate     string
	Timestamp uint32
}

type IAmACamera struct {
	Road  uint16
	Mile  uint16
	Limit uint16
}

type ParseFunc[T any] func(b []byte) (T, []byte)

func parsePlate(b []byte) (*Plate, []byte) {
	plate, b := parseString(b)
	timestamp, b := parseUint32(b)

	return &Plate{
		Plate:     plate,
		Timestamp: timestamp,
	}, b
}

func parseIAmACamera(b []byte) (*IAmACamera, []byte) {
	road, b := parseUint16(b)
	mile, b := parseUint16(b)
	limit, b := parseUint16(b)

	return &IAmACamera{
		Road:  road,
		Mile:  mile,
		Limit: limit,
	}, b
}

// Consumes tokens from b to produce a string
// Returns number of bytes consumed and the final string
func parseString(b []byte) (string, []byte) {
	len := int(b[0])
	str := string(b[1 : len+1])

	return str, b[len+1:]
}

// parse uint16
func parseUint8(b []byte) (uint8, []byte) {
	return uint8(b[0]), b[1:]
}

// parse uint16
func parseUint16(b []byte) (uint16, []byte) {
	return binary.BigEndian.Uint16(b), b[2:]
}

// parse uint32
func parseUint32(b []byte) (uint32, []byte) {
	return binary.BigEndian.Uint32(b), b[4:]
}
