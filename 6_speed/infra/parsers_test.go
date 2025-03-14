package infra_test

import (
	"protohackers/6_speed/infra"
	"slices"
	"testing"
)

type ParsingCases[T any] struct {
	in       []byte
	expected T
	newLen   int
	fn       infra.ParseFunc[T]
	eq       func(T, T) bool
}

func runParsingCases[T any](t *testing.T, cases []ParsingCases[T]) {
	for _, c := range cases {
		out := c.fn(c.in)
		if !out.Ok {
			t.Fatalf("parsing failed")
		}
		if !c.eq(out.Value, c.expected) {
			t.Fatalf("wrong parsing output. expected %v got %v", c.expected, out)
		}
		if len(out.Next) != c.newLen {
			t.Fatalf("wrong new input length. expected %v got %v", c.newLen, len(out.Next))
		}
	}
}

func TestParser(t *testing.T) {
	pleq := func(p1, p2 *infra.Plate) bool {
		return p1.Plate == p2.Plate && p1.Timestamp == p2.Timestamp
	}

	plateCases := []ParsingCases[*infra.Plate]{
		{
			in: []byte{0x20, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8},
			expected: &infra.Plate{
				Plate:     "UN1X",
				Timestamp: 1000,
			},
			newLen: 0,
			fn:     infra.ParsePlate,
			eq:     pleq,
		},
	}

	cameq := func(ia1, ia2 *infra.IAmACamera) bool {
		return ia1.Limit == ia2.Limit && ia1.Mile == ia2.Mile && ia1.Road == ia2.Road
	}
	cameraCases := []ParsingCases[*infra.IAmACamera]{
		{
			in: []byte{0x80, 0x00, 0x42, 0x00, 0x64, 0x00, 0x3c},
			expected: &infra.IAmACamera{
				Road:  66,
				Mile:  100,
				Limit: 60,
			},
			newLen: 0,
			fn:     infra.ParseIAmACamera,
			eq:     cameq,
		},
	}

	dispeq := func(ia1, ia2 *infra.IAmADispatcher) bool {
		return slices.Equal(ia1.Roads, ia2.Roads)
	}
	dispatchCases := []ParsingCases[*infra.IAmADispatcher]{
		{
			in: []byte{0x81, 0x03, 0x00, 0x42, 0x01, 0x70, 0x13, 0x88},
			expected: &infra.IAmADispatcher{
				Roads: []uint16{66, 368, 5000},
			},
			newLen: 0,
			fn:     infra.ParseIAmADispatcher,
			eq:     dispeq,
		},
	}

	hbeq := func(wh1, wh2 *infra.WantHeartbeat) bool { return wh1.Interval == wh2.Interval }
	hbcases := []ParsingCases[*infra.WantHeartbeat]{
		{
			in: []byte{0x40, 0x00, 0x00, 0x00, 0x10},
			expected: &infra.WantHeartbeat{
				Interval: 16,
			},
			newLen: 0,
			fn:     infra.ParseWantHeartbeat,
			eq:     hbeq,
		},
	}

	t.Run("plate cases", func(t *testing.T) {
		runParsingCases(t, plateCases)
	})
	t.Run("camera cases", func(t *testing.T) {
		runParsingCases(t, cameraCases)
	})
	t.Run("dispatch cases", func(t *testing.T) {
		runParsingCases(t, dispatchCases)
	})
	t.Run("hb cases", func(t *testing.T) {
		runParsingCases(t, hbcases)
	})
}
