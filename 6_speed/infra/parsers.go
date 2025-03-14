package infra

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ParseFunc[T any] func(b []byte) ParseResult[T]
type ParseResult[T any] struct {
	Value T
	Next  []byte
	Ok    bool
}

func ParseMessages(r io.Reader) chan any {
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

			plate := ParsePlate(curr)
			if plate.Ok {
				curr = plate.Next
				ch <- plate.Value
			}

			hb := ParseWantHeartbeat(curr)
			if hb.Ok {
				curr = hb.Next
				ch <- hb
			}
		}
	}()

	return ch
}

func ParsePlate(b []byte) ParseResult[*Plate] {
	var ret ParseResult[*Plate]

	typeHex := ParseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x20 {
		return ret
	}

	plate := ParseString(typeHex.Next)
	if !plate.Ok {
		return ret
	}

	timestamp := ParseUint32(plate.Next)
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

func ParseIAmACamera(b []byte) ParseResult[*IAmACamera] {
	var ret ParseResult[*IAmACamera]

	typeHex := ParseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x80 {
		return ret
	}

	road := ParseUint16(typeHex.Next)
	if !road.Ok {
		return ret
	}
	mile := ParseUint16(road.Next)
	if !mile.Ok {
		return ret
	}
	limit := ParseUint16(mile.Next)
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

	typeHex := ParseUint8(b)
	if !typeHex.Ok || typeHex.Value != 0x81 {
		return ret
	}

	numroads := ParseUint8(typeHex.Next)
	if !numroads.Ok {
		return ret
	}

	var i uint8
	next := numroads.Next
	roads := make([]uint16, 0)

	for i = 0; i < numroads.Value; i++ {
		road := ParseUint16(next)
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

func ParseWantHeartbeat(b []byte) ParseResult[*WantHeartbeat] {
	var ret ParseResult[*WantHeartbeat]

	typeHex := ParseUint8(b)
	if typeHex.Ok && typeHex.Value != 0x40 {
		return ret
	}

	out := ParseUint32(typeHex.Next)
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
func ParseString(b []byte) ParseResult[string] {
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
func ParseUint8(b []byte) ParseResult[uint8] {
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
func ParseUint16(b []byte) ParseResult[uint16] {
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
func ParseUint32(b []byte) ParseResult[uint32] {
	var ret ParseResult[uint32]
	if len(b) < 4 {
		return ret
	}
	ret.Ok = true
	ret.Value = binary.BigEndian.Uint32(b)
	ret.Next = b[4:]
	return ret
}
