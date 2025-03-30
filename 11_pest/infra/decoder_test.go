package infra

import (
	"testing"
)

func TestU32(t *testing.T) {
	p := parseUint32([]byte{0xa6, 0xa9, 0xb5, 0x67})
	exp := 2796139879
	if int(p.Value) != exp {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestString(t *testing.T) {
	p := parseString([]byte{0x00, 0x00, 0x00, 0x03, 0x66, 0x6f, 0x6f})
	exp := "foo"
	if p.Value != exp {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}
