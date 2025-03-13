package main

import (
	"testing"
)

func TestParseSpace(t *testing.T) {
	type spaceCases struct {
		in          []byte
		expected    string
		parsedBytes int
	}

	cases := []spaceCases{
		{
			in:          []byte{0x03, 0x66, 0x6f, 0x6f, 0x00, 0x05},
			expected:    "foo",
			parsedBytes: 4,
		},
		{
			in:          []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected:    "Elbereth",
			parsedBytes: 9,
		},
	}
	for _, c := range cases {
		out, n := parseString(c.in)

		if out != c.expected {
			t.Fatalf("failed to parse. expected %v got %v", c.expected, out)
		}
		if n != c.parsedBytes {
			t.Fatalf("consumed different bytes than expected. expected %v got %v", c.parsedBytes, n)
		}
	}

}

func TestApplyParser(t *testing.T) {
	type applyCases struct {
		in       []byte
		expected int
		newLen   int
	}

	cases := []applyCases{
		{
			in:       []byte{0x03, 0x66, 0x6f, 0x6f, 0x00, 0x05},
			expected: 0x03,
			newLen:   5,
		},
		{
			in:       []byte{0x01},
			expected: 0x01,
			newLen:   0,
		},
		{
			in:       []byte{0x08, 0x45, 0x6C, 0x62, 0x65, 0x72, 0x65, 0x74, 0x68, 0x01, 0x02},
			expected: 0x08,
			newLen:   10,
		},
	}
	for _, c := range cases {
		out, new_in := applyParser(parseUint8, c.in)
		if out != uint8(c.expected) {
			t.Fatalf("wrong parsing output. expected %v got %v", c.expected, out)
		}
		if len(new_in) != c.newLen {
			t.Fatalf("wrong new input length. expected %v got %v", c.newLen, len(new_in))
		}
	}

}
