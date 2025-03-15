package infra

import (
	"context"
	"encoding/binary"
	"io"
)

type ParseFunc[T any] func(b []byte) ParseResult[T]
type ParseResult[T any] struct {
	Value T
	Next  []byte
	Ok    bool
}

func ParseMessages(r io.Reader, cancel context.CancelFunc) chan any {
	ch := make(chan any)

	go func() {
		defer close(ch)
		if cancel != nil {
			defer cancel()
		}

		var curr []byte

		for {
			buff := make([]byte, 1024)
			n, err := r.Read(buff)
			if err != nil {
				return
			}
			curr = append(curr, buff[:n]...)

			for {
				newcurr := executeParse(ch, curr)

				noMoreToParse := len(newcurr) == 0
				nothingParsed := len(newcurr) == len(curr)
				curr = newcurr

				if noMoreToParse || nothingParsed {
					break
				}
			}
		}
	}()

	return ch
}

func executeParse(ch chan any, curr []byte) []byte {
	spe_err := ParseError(curr)
	if spe_err.Ok {
		curr = spe_err.Next
		ch <- spe_err.Value
	}

	plate := ParsePlate(curr)
	if plate.Ok {
		curr = plate.Next
		ch <- plate.Value
	}

	tick := ParseTicket(curr)
	if tick.Ok {
		curr = tick.Next
		ch <- tick.Value
	}

	whb := ParseWantHeartbeat(curr)
	if whb.Ok {
		curr = whb.Next
		ch <- whb.Value
	}

	hb := ParseHeartbeat(curr)
	if hb.Ok {
		curr = hb.Next
		ch <- hb.Value
	}

	cam := ParseIAmACamera(curr)
	if cam.Ok {
		curr = cam.Next
		ch <- cam.Value
	}

	disp := ParseIAmADispatcher(curr)
	if disp.Ok {
		curr = disp.Next
		ch <- disp.Value
	}

	return curr
}

func ParseError(b []byte) ParseResult[*SpeedError] {
	var ret ParseResult[*SpeedError]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x10 {
		return ret
	}

	msg := parseString(typeHex.Next)
	if !msg.Ok {
		return ret
	}

	ret.Ok = true
	ret.Next = msg.Next
	ret.Value = &SpeedError{
		Msg: msg.Value,
	}
	return ret
}

func ParsePlate(b []byte) ParseResult[*Plate] {
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

func ParseTicket(b []byte) ParseResult[*Ticket] {
	var ret ParseResult[*Ticket]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x21 {
		return ret
	}

	plate := parseString(typeHex.Next)
	if !plate.Ok {
		return ret
	}

	road := parseUint16(plate.Next)
	if !road.Ok {
		return ret
	}

	mile1 := parseUint16(road.Next)
	if !mile1.Ok {
		return ret
	}

	timestamp1 := parseUint32(mile1.Next)
	if !timestamp1.Ok {
		return ret
	}

	mile2 := parseUint16(timestamp1.Next)
	if !mile2.Ok {
		return ret
	}

	timestamp2 := parseUint32(mile2.Next)
	if !timestamp2.Ok {
		return ret
	}

	speed := parseUint16(timestamp2.Next)
	if !speed.Ok {
		return ret
	}

	ret.Ok = true
	ret.Next = speed.Next
	ret.Value = &Ticket{
		Plate:      plate.Value,
		Road:       road.Value,
		Mile1:      mile1.Value,
		Timestamp1: timestamp1.Value,
		Mile2:      mile2.Value,
		Timestamp2: timestamp2.Value,
		Speed:      speed.Value,
	}
	return ret
}

func ParseWantHeartbeat(b []byte) ParseResult[*WantHeartbeat] {
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

func ParseHeartbeat(b []byte) ParseResult[Heartbeat] {
	var ret ParseResult[Heartbeat]

	typeHex := parseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x41 {
		return ret
	}

	ret.Ok = true
	ret.Next = typeHex.Next

	return ret
}

func ParseIAmACamera(b []byte) ParseResult[*IAmACamera] {
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

func ParseIAmADispatcher(b []byte) ParseResult[*IAmADispatcher] {
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
