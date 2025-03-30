package infra

import (
	"reflect"
	"testing"
)

func TestU32(t *testing.T) {
	p := parseUint32([]byte{0xa6, 0xa9, 0xb5, 0x67})
	exp := uint32(2796139879)
	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestString(t *testing.T) {
	p := parseString([]byte{0x00, 0x00, 0x00, 0x03, 0x66, 0x6f, 0x6f})
	exp := "foo"
	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}

func TestArray(t *testing.T) {
	parser := parseArray(parseUint32)

	p := parser([]byte{0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x00, 0x00, 0x0b})
	exp := []uint32{10, 11}

	if !reflect.DeepEqual(p.Value, exp) {
		t.Fatalf("wrong parse result. exp %v got %v", exp, p.Value)
	}
}
