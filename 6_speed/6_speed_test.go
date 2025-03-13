package main

import (
	"testing"
)

type ParsingCases[T any] struct {
	in       []byte
	expected T
	newLen   int
	fn       ParseFunc[T]
}

func runParsingCases[T comparable](t *testing.T, cases []ParsingCases[T]) {
	for _, c := range cases {
		out, new_in := c.fn(c.in)
		if out != c.expected {
			t.Fatalf("wrong parsing output. expected %v got %v", c.expected, out)
		}
		if len(new_in) != c.newLen {
			t.Fatalf("wrong new input length. expected %v got %v", c.newLen, len(new_in))
		}
	}
}

func TestParser(t *testing.T) {
	uint8Cases := []ParsingCases[uint8]{
		{
			in:       []byte{0x03, 0x66, 0x6f, 0x6f, 0x00, 0x05},
			expected: 0x03,
			newLen:   5,
			fn:       parseUint8,
		},
		{
			in:       []byte{0x01},
			expected: 0x01,
			newLen:   0,
			fn:       parseUint8,
		},
		{
			in:       []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected: 0x08,
			newLen:   10,
			fn:       parseUint8,
		},
	}

	stringCases := []ParsingCases[string]{
		{
			in:       []byte{0x03, 0x66, 0x6f, 0x6f, 0x00},
			expected: "foo",
			newLen:   1,
			fn:       parseString,
		},
		{
			in:       []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected: "Elbereth",
			newLen:   2,
			fn:       parseString,
		},
	}

	t.Run("uint8 cases", func(t *testing.T) {
		runParsingCases(t, uint8Cases)
	})

	t.Run("string cases", func(t *testing.T) {
		runParsingCases(t, stringCases)
	})

}

// i love generics
