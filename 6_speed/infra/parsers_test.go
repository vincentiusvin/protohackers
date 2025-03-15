package infra_test

import (
	"bytes"
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

func TestIndividualParser(t *testing.T) {
	t.Run("error cases", func(t *testing.T) {
		erreq := func(p1, p2 *infra.SpeedError) bool {
			return p1.Msg == p2.Msg
		}
		errCases := []ParsingCases[*infra.SpeedError]{
			{
				in: []byte{0x10, 0x03, 0x62, 0x61, 0x64},
				expected: &infra.SpeedError{
					Msg: "bad",
				},
				newLen: 0,
				fn:     infra.ParseError,
				eq:     erreq,
			},
			{
				in: []byte{0x10, 0x03, 0x62, 0x61, 0x64, 0x20},
				expected: &infra.SpeedError{
					Msg: "bad",
				},
				newLen: 1,
				fn:     infra.ParseError,
				eq:     erreq,
			},
		}
		runParsingCases(t, errCases)
	})

	t.Run("plate cases", func(t *testing.T) {
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
		runParsingCases(t, plateCases)
	})

	t.Run("ticket cases", func(t *testing.T) {
		tickeq := func(p1, p2 *infra.Ticket) bool {
			return p1.Mile1 == p2.Mile1 &&
				p1.Mile2 == p2.Mile2 &&
				p1.Plate == p2.Plate &&
				p1.Road == p2.Road &&
				p1.Speed == p2.Speed &&
				p1.Timestamp1 == p2.Timestamp1 &&
				p1.Timestamp2 == p2.Timestamp2
		}
		tickCases := []ParsingCases[*infra.Ticket]{
			{
				in: []byte{
					0x21,
					0x04, 0x55, 0x4e, 0x31, 0x58,
					0x00, 0x42,
					0x00, 0x64,
					0x00, 0x01, 0xe2, 0x40,
					0x00, 0x6e,
					0x00, 0x01, 0xe3, 0xa8,
					0x27, 0x10,
				},
				expected: &infra.Ticket{
					Plate:      "UN1X",
					Road:       66,
					Mile1:      100,
					Timestamp1: 123456,
					Mile2:      110,
					Timestamp2: 123816,
					Speed:      10000,
				},
				newLen: 0,
				fn:     infra.ParseTicket,
				eq:     tickeq,
			},
		}
		runParsingCases(t, tickCases)
	})

	t.Run("hb cases", func(t *testing.T) {
		hbeq := func(wh1, wh2 infra.Heartbeat) bool { return true }
		hbcases := []ParsingCases[infra.Heartbeat]{
			{
				in:       []byte{0x41},
				expected: infra.Heartbeat{},
				newLen:   0,
				fn:       infra.ParseHeartbeat,
				eq:       hbeq,
			},
		}
		runParsingCases(t, hbcases)
	})

	t.Run("want hb cases", func(t *testing.T) {
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
		runParsingCases(t, hbcases)
	})

	t.Run("camera cases", func(t *testing.T) {
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
		runParsingCases(t, cameraCases)
	})
	t.Run("dispatch cases", func(t *testing.T) {
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
		runParsingCases(t, dispatchCases)
	})
}

func TestCombinedParser(t *testing.T) {
	in := []byte{
		0x20, 0x04, 0x55, 0x4e, 0x31, 0x58, 0x00, 0x00, 0x03, 0xe8,
		0x20, 0x04, 0x55, 0x4e, 0x32, 0x58, 0x00, 0x00, 0x03, 0xe8,
	}
	in_r := new(bytes.Buffer)
	in_r.Write(in)

	c := infra.ParseMessages(in_r, nil)

	pleq := func(p1, p2 *infra.Plate) bool {
		return p1.Plate == p2.Plate && p1.Timestamp == p2.Timestamp
	}

	out1 := (<-c).(*infra.Plate)
	exp1 := &infra.Plate{
		Plate:     "UN1X",
		Timestamp: 1000,
	}
	if !pleq(out1, exp1) {
		t.Fatalf("wrong output. expect %v got %v", exp1, out1)
	}

	out2 := (<-c).(*infra.Plate)
	exp2 := &infra.Plate{
		Plate:     "UN2X",
		Timestamp: 1000,
	}
	if !pleq(out2, exp2) {
		t.Fatalf("wrong output. expect %v got %v", exp2, out2)
	}

}
