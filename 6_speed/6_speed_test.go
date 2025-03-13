package main

import (
	"slices"
	"testing"
)

type ParsingCases[T any] struct {
	in       []byte
	expected T
	newLen   int
	fn       ParseFunc[T]
	eq       func(T, T) bool
}

func runParsingCases[T any](t *testing.T, cases []ParsingCases[T]) {
	for _, c := range cases {
		out, new_in := c.fn(c.in)
		if !c.eq(out, c.expected) {
			t.Fatalf("wrong parsing output. expected %v got %v", c.expected, out)
		}
		if len(new_in) != c.newLen {
			t.Fatalf("wrong new input length. expected %v got %v", c.newLen, len(new_in))
		}
	}
}

func TestParser(t *testing.T) {
	uint8eq := func(u1, u2 uint8) bool { return u1 == u2 }
	uint8Cases := []ParsingCases[uint8]{
		{
			in:       []byte{0x03, 0x66, 0x6f, 0x6f, 0x00, 0x05},
			expected: 0x03,
			newLen:   5,
			fn:       parseUint8,
			eq:       uint8eq,
		},
		{
			in:       []byte{0x01},
			expected: 0x01,
			newLen:   0,
			fn:       parseUint8,
			eq:       uint8eq,
		},
		{
			in:       []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected: 0x08,
			newLen:   10,
			fn:       parseUint8,
			eq:       uint8eq,
		},
	}

	streq := func(s1, s2 string) bool { return s1 == s2 }
	stringCases := []ParsingCases[string]{
		{
			in:       []byte{0x03, 0x66, 0x6f, 0x6f, 0x00},
			expected: "foo",
			newLen:   1,
			fn:       parseString,
			eq:       streq,
		},
		{
			in:       []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected: "Elbereth",
			newLen:   2,
			fn:       parseString,
			eq:       streq,
		},
	}

	pleq := func(p1, p2 *Plate) bool {
		return p1.Plate == p2.Plate && p1.Timestamp == p2.Timestamp
	}

	plateCases := []ParsingCases[*Plate]{
		{
			in: []byte{0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8},
			expected: &Plate{
				Plate:     "UN1X",
				Timestamp: 1000,
			},
			newLen: 0,
			fn:     parsePlate,
			eq:     pleq,
		},
	}

	cameq := func(ia1, ia2 *IAmACamera) bool {
		return ia1.Limit == ia2.Limit && ia1.Mile == ia2.Mile && ia1.Road == ia2.Road
	}
	cameraCases := []ParsingCases[*IAmACamera]{
		{
			in: []byte{0x00, 0x42, 0x00, 0x64, 0x00, 0x3c},
			expected: &IAmACamera{
				Road:  66,
				Mile:  100,
				Limit: 60,
			},
			newLen: 0,
			fn:     parseIAmACamera,
			eq:     cameq,
		},
	}

	dispeq := func(ia1, ia2 *IAmADispatcher) bool {
		return slices.Equal[[]uint16](ia1.Roads, ia2.Roads)
	}
	dispatchCases := []ParsingCases[*IAmADispatcher]{
		{
			in: []byte{0x03, 0x00, 0x42, 0x01, 0x70, 0x13, 0x88},
			expected: &IAmADispatcher{
				Roads: []uint16{66, 368, 5000},
			},
			newLen: 0,
			fn:     parseIAmADispatcher,
			eq:     dispeq,
		},
	}

	t.Run("uint8 cases", func(t *testing.T) {
		runParsingCases(t, uint8Cases)
	})

	t.Run("string cases", func(t *testing.T) {
		runParsingCases(t, stringCases)
	})

	t.Run("plate cases", func(t *testing.T) {
		runParsingCases(t, plateCases)
	})
	t.Run("camera cases", func(t *testing.T) {
		runParsingCases(t, cameraCases)
	})
	t.Run("dispatch cases", func(t *testing.T) {
		runParsingCases(t, dispatchCases)
	})

}

// i love generics
