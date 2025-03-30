package infra

import (
	"encoding/binary"
)

type ParseFunc[T any] func(b []byte) ParseResult[T]
type ParseResult[T any] struct {
	Value T
	Next  []byte
	Ok    bool
}

type Hello struct {
	Protocol string
	Version  uint32
}

func parseHello(b []byte) (ret ParseResult[Hello]) {
	prefix := parseUint8(b)
	if !prefix.Ok || prefix.Value != uint8(0x50) {
		return ret
	}
	msgLen := parseUint32(prefix.Next)
	if !msgLen.Ok {
		return ret
	}
	protocol := parseString(msgLen.Next)
	if !protocol.Ok {
		return ret
	}
	version := parseUint32(protocol.Next)
	if !version.Ok {
		return ret
	}
	checksum := parseUint8(version.Next)
	if !checksum.Ok {
		return ret
	}

	expectedMsgLen := int(msgLen.Value)
	actualMsgLen := len(b) - len(checksum.Next)
	if expectedMsgLen != actualMsgLen {
		return ret
	}

	var sum uint8
	for i := 0; i < actualMsgLen; i++ {
		sum += b[i]
	}
	if sum != 0 {
		return ret
	}

	ret.Ok = true
	ret.Value = Hello{
		Protocol: protocol.Value,
		Version:  version.Value,
	}
	ret.Next = version.Next
	return ret
}

func parseUint8(b []byte) (ret ParseResult[uint8]) {
	if len(b) < 1 {
		return
	}
	ret.Ok = true
	ret.Value = uint8(b[0])
	ret.Next = b[1:]
	return ret
}

func parseUint32(b []byte) (ret ParseResult[uint32]) {
	if len(b) < 4 {
		return
	}
	ret.Ok = true
	ret.Value = binary.BigEndian.Uint32(b)
	ret.Next = b[4:]
	return ret
}

// Consumes tokens from b to produce a string
// Returns number of bytes consumed and the final string
func parseString(b []byte) (ret ParseResult[string]) {
	lenParse := parseUint32(b)
	if !lenParse.Ok {
		return
	}

	b = lenParse.Next
	lenVal := int(lenParse.Value)

	if len(b) < lenVal {
		return ret
	}

	str := string(b[:lenVal])

	ret.Ok = true
	ret.Value = str
	ret.Next = b[lenVal:]

	return ret
}

// parser combinator :)
func parseArray[T any](fn ParseFunc[T]) ParseFunc[[]T] {
	return func(b []byte) (ret ParseResult[[]T]) {
		lenParse := parseUint32(b)
		if !lenParse.Ok {
			return ret
		}

		b = lenParse.Next
		lenVal := int(lenParse.Value)

		acc := make([]T, lenVal)
		for i := 0; i < lenVal; i++ {
			curr := fn(b)
			if !curr.Ok {
				return ret
			}
			acc[i] = curr.Value
			b = curr.Next
		}

		ret.Ok = true
		ret.Value = acc
		ret.Next = b

		return ret
	}
}
