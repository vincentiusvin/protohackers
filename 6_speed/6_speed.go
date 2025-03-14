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

type ParseFunc[T any] func(b []byte) ParseResult[T]
type ParseResult[T any] struct {
	Value T
	Next  []byte
	Ok    bool
}

type ThenResult[T any, U any] struct {
	Value1 T
	Value2 U
}

func then[T any, U any](fn1 ParseFunc[T], fn2 ParseFunc[U]) ParseFunc[ThenResult[T, U]] {
	return func(b []byte) ParseResult[ThenResult[T, U]] {
		var ret ParseResult[ThenResult[T, U]]
		val1 := fn1(b)
		if !val1.Ok {
			return ret
		}

		val2 := fn2(val1.Next)
		if !val2.Ok {
			return ret
		}

		ret.Ok = true
		ret.Value = ThenResult[T, U]{
			Value1: val1.Value,
			Value2: val2.Value,
		}
		ret.Next = val2.Next

		return ret
	}
}

func parsePlate(b []byte) ParseResult[*Plate] {
	var ret ParseResult[*Plate]

	iamgoingtohell := then(then(parseUint8, parseString), parseUint32)(b)

	hexCode := iamgoingtohell.Value.Value1.Value1
	if hexCode != 0x20 {
		return ret
	}

	plate := iamgoingtohell.Value.Value1.Value2
	timestamp := iamgoingtohell.Value.Value2

	ret.Ok = true
	ret.Value = &Plate{
		Plate:     plate,
		Timestamp: timestamp,
	}
	ret.Next = iamgoingtohell.Next
	return ret
}

func parseIAmACamera(b []byte) ParseResult[*IAmACamera] {
	var ret ParseResult[*IAmACamera]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x80 {
		return ret
	}

	road := parseUint16(typeHex.Next)
	if !road.Ok {
		return ret
	}
	mile := parseUint16(road.Next)
	if !mile.Ok {
		return ret
	}
	limit := parseUint16(mile.Next)
	if !limit.Ok {
		return ret
	}

	ret.Ok = true
	ret.Value = &IAmACamera{
		Road:  road.Value,
		Mile:  mile.Value,
		Limit: limit.Value,
	}
	ret.Next = limit.Next
	return ret
}

func parseIAmADispatcher(b []byte) ParseResult[*IAmADispatcher] {
	var ret ParseResult[*IAmADispatcher]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x81 {
		return ret
	}

	numroads := parseUint8(typeHex.Next)
	if !numroads.Ok {
		return ret
	}

	var i uint8
	next := numroads.Next
	roads := make([]uint16, 0)

	for i = 0; i < numroads.Value; i++ {
		road := parseUint16(next)
		if !road.Ok {
			return ret
		}
		next = road.Next
		roads = append(roads, road.Value)
	}

	ret.Ok = true
	ret.Value = &IAmADispatcher{
		Roads: roads,
	}
	ret.Next = next

	return ret
}

func parseWantHeartbeat(b []byte) ParseResult[*WantHeartbeat] {
	var ret ParseResult[*WantHeartbeat]

	typeHex := parseUint8(b)
	if typeHex.Ok && typeHex.Value != 0x40 {
		return ret
	}

	out := parseUint32(typeHex.Next)
	if !out.Ok {
		return ret
	}

	ret.Ok = true
	ret.Value = &WantHeartbeat{
		Interval: out.Value,
	}
	ret.Next = out.Next

	return ret
}

// Consumes tokens from b to produce a string
// Returns number of bytes consumed and the final string
func parseString(b []byte) ParseResult[string] {
	var ret ParseResult[string]

	if len(b) < 1 {
		return ret
	}

	strlen := int(b[0])
	if len(b) < strlen+1 {
		return ret
	}

	str := string(b[1 : strlen+1])

	ret.Ok = true
	ret.Value = str
	ret.Next = b[strlen+1:]

	return ret
}

// parse uint8
func parseUint8(b []byte) ParseResult[uint8] {
	var ret ParseResult[uint8]
	if len(b) < 1 {
		return ret
	}
	ret.Ok = true
	ret.Value = uint8(b[0])
	ret.Next = b[1:]
	return ret
}

// parse uint16
func parseUint16(b []byte) ParseResult[uint16] {
	var ret ParseResult[uint16]
	if len(b) < 2 {
		return ret
	}
	ret.Ok = true
	ret.Value = binary.BigEndian.Uint16(b)
	ret.Next = b[2:]
	return ret
}

// parse uint32
func parseUint32(b []byte) ParseResult[uint32] {
	var ret ParseResult[uint32]
	if len(b) < 4 {
		return ret
	}
	ret.Ok = true
	ret.Value = binary.BigEndian.Uint32(b)
	ret.Next = b[4:]
	return ret
}

type Ticket struct {
	Plate      string
	Road       uint16
	Mile1      uint16
	Timestamp1 uint32
	Mile2      uint16
	Timestamp2 uint32
	Speed      uint16
}

func encodeTicket(t *Ticket) []byte {
	ret := append(
		[]byte{0x21},
		encodeString(t.Plate)...,
	)
	ret = binary.BigEndian.AppendUint16(ret, t.Road)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile1)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp1)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile2)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp2)
	ret = binary.BigEndian.AppendUint16(ret, t.Speed)
	return ret
}

func encodeHeartbeat() []byte {
	return []byte{0x41}
}

func encodeError(err string) []byte {
	return append([]byte{0x10}, encodeString(err)...)
}

func encodeString(s string) []byte {
	l := byte(uint8(len(s)))
	return append([]byte{l}, []byte(s)...)
}
