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

type IAmADispatcher struct {
	Roads []uint16
}

type WantHeartbeat struct {
	Interval uint32
}

type ParseFunc[T any] func(b []byte) (T, []byte)

func parseInfer(b []byte) (any, []byte) {
	if len(b) == 0 {
		return nil, b
	}
	typeHex, b := parseUint8(b)
	switch typeHex {
	case 0x20:
		return parsePlate(b)
	case 0x40:
		return parseWantHeartbeat(b)
	case 0x80:
		return parseIAmACamera(b)
	case 0x81:
		return parseIAmADispatcher(b)
	}
	return nil, b
}

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

func parseIAmADispatcher(b []byte) (*IAmADispatcher, []byte) {
	numroads, b := parseUint8(b)
	var i uint8
	roads := make([]uint16, 0)

	for i = 0; i < numroads; i++ {
		road, new_b := parseUint16(b)
		b = new_b
		roads = append(roads, road)

	}

	return &IAmADispatcher{
		Roads: roads,
	}, b
}

func parseWantHeartbeat(b []byte) (*WantHeartbeat, []byte) {
	hb, b := parseUint32(b)

	return &WantHeartbeat{
		Interval: hb,
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
