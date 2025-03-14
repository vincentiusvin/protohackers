package infra_test

import (
	"bytes"
	"protohackers/6_speed/infra"
	"testing"
)

func TestEncode(t *testing.T) {
	ticket := &infra.Ticket{
		Plate:      "UN1X",
		Road:       66,
		Mile1:      100,
		Timestamp1: 123456,
		Mile2:      110,
		Timestamp2: 123816,
		Speed:      10000,
	}

	out := infra.EncodeTicket(ticket)

	expected := []byte{
		0x21,
		0x04, 0x55, 0x4e, 0x31, 0x58,
		0x00, 0x42,
		0x00, 0x64,
		0x00, 0x01, 0xe2, 0x40,
		0x00, 0x6e,
		0x00, 0x01, 0xe3, 0xa8,
		0x27, 0x10,
	}

	if !bytes.Equal(expected, out) {
		t.Fatalf("failed to serialize %v", ticket)
	}
}
