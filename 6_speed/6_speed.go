package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
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

type ClientType int

const (
	None ClientType = iota
	Camera
	Dispatcher
)

func handleConnection(c net.Conn) {
	defer c.Close()

	ch := parseMessages(c)
	send := func(b []byte) error {
		_, err := c.Write(b)
		return err
	}
	handleConnectionLogic(ch, send)
}

// With the connection abstracted away
func handleConnectionLogic(c chan any, send func([]byte) error) {
	clientType := None
	spawnedHb := false

	for msg := range c {
		switch v := msg.(type) {
		case *IAmACamera:
			if clientType != None {
				send(encodeError("already another type"))
				continue
			}

			clientType = Camera
			log.Println("new camera", v)
		case *IAmADispatcher:
			if clientType != None {
				send(encodeError("already another type"))
				continue
			}

			clientType = Dispatcher
			log.Println("new dispatcher", v)
		case *WantHeartbeat:
			if spawnedHb {
				send(encodeError("already sending heartbeat"))
				continue
			}
			go func() {
				log.Println("starting heartbeat every", v.Interval, "ds")
				for {
					err := send(encodeHeartbeat())
					if err != nil {
						break
					}
					deciseconds := time.Second / 10
					time.Sleep(time.Duration(v.Interval) * deciseconds)
				}
				log.Println("stopped heartbeat every", v.Interval, "ds")
			}()
			spawnedHb = true
		case *Plate:
			if clientType != Camera {
				send(encodeError("you are not a camera"))
				continue
			}
		}
	}
}

func parseMessages(r io.Reader) chan any {
	buff := new(bytes.Buffer)
	buff.ReadFrom(r)

	ch := make(chan any)

	go func() {
		defer close(ch)
		var curr []byte
		for {
			b, err := buff.ReadByte()
			if err != nil {
				break
			}
			curr = append(curr, b)

			cam := parseIAmACamera(curr)
			if cam.Ok {
				curr = cam.Next
				ch <- cam.Value
			}

			disp := parseIAmADispatcher(curr)
			if disp.Ok {
				curr = disp.Next
				ch <- disp.Value
			}

			plate := parsePlate(curr)
			if plate.Ok {
				curr = plate.Next
				ch <- plate.Value
			}

			hb := parseWantHeartbeat(curr)
			if hb.Ok {
				curr = hb.Next
				ch <- hb
			}
		}
	}()

	return ch
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

func parsePlate(b []byte) ParseResult[*Plate] {
	var ret ParseResult[*Plate]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x20 {
		return ret
	}

	plate := parseString(typeHex.Next)
	if !plate.Ok {
		return ret
	}

	timestamp := parseUint32(plate.Next)
	if !timestamp.Ok {
		return ret
	}

	ret.Ok = true
	ret.Value = &Plate{
		Plate:     plate.Value,
		Timestamp: timestamp.Value,
	}
	ret.Next = timestamp.Next
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
