package main

import (
	"testing"
)

func TestParseInsert(t *testing.T) {
	in := []byte{0x49, 0x00, 0x00, 0x30, 0x39, 0x00, 0x00, 0x00, 0x65}
	exp_price := 101
	exp_timestamp := 12345

	insert, _ := parsePacket(in)
	if insert == nil {
		t.Fatalf("Failed to parse packet")
	}
	if insert.Price != exp_price {
		t.Fatalf("wrong price. got %v expected %v", insert.Price, exp_price)
	}
	if insert.Timestamp != exp_timestamp {
		t.Fatalf("wrong timestamp. got %v expected %v", insert.Timestamp, exp_timestamp)
	}
}

func TestParseQuery(t *testing.T) {
	in := []byte{0x51, 0x00, 0x00, 0x30, 0x39, 0xFF, 0xFF, 0xFF, 0xFF}
	exp_mintime := 12345
	exp_maxtime := -1

	_, query := parsePacket(in)
	if query == nil {
		t.Fatalf("Failed to parse packet")
	}
	if query.Mintime != exp_mintime {
		t.Fatalf("wrong mintime. got %v expected %v", query.Mintime, exp_mintime)
	}
	if query.Maxtime != exp_maxtime {
		t.Fatalf("wrong maxtime. got %v expected %v", query.Maxtime, exp_maxtime)
	}
}
