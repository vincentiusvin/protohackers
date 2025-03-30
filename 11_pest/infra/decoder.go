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
			return
		}

		b = lenParse.Next
		lenVal := int(lenParse.Value)

		acc := make([]T, lenVal)
		for i := 0; i < lenVal; i++ {
			curr := fn(b)
			if !curr.Ok {
				return
			}
			acc[i] = curr.Value
			b = curr.Next
		}

		ret.Ok = true
		ret.Value = acc
		ret.Next = b

		return
	}
}
