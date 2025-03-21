package cipher_test

import (
	"protohackers/8_cipher/cipher"
	"reflect"
	"testing"
)

func TestCipher(t *testing.T) {
	type cipherCases struct {
		ciph []byte
		in   []byte
		exp  []byte
	}
	cases := []cipherCases{
		{
			ciph: []byte{0x01},
			in:   []byte{0x69, 0x64, 0x6d, 0x6d, 0x6e},
			exp:  []byte{0x96, 0x26, 0xb6, 0xb6, 0x76},
		},
		{
			ciph: []byte{0x02, 0x01},
			in:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			exp:  []byte{0x69, 0x64, 0x6d, 0x6d, 0x6e},
		},
		{
			ciph: []byte{0x03},
			in:   []byte{0xFF, 0xFF, 0xFF},
			exp:  []byte{0xFF, 0xFE, 0xFD},
		},
		{
			ciph: []byte{0x04, 0x02},
			in:   []byte{0x00, 0xFF, 0xF0},
			exp:  []byte{0x02, 0x01, 0xF2},
		},
		{
			ciph: []byte{0x05},
			in:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			exp:  []byte{0x68, 0x66, 0x6e, 0x6f, 0x73},
		},
	}

	for _, c := range cases {
		ciph := cipher.ParseCipher(c.ciph)
		out := ciph.Encode(c.in)
		if !reflect.DeepEqual(out, c.exp) {
			t.Fatalf("cipher wrong. exp: %v, got: %v", c.exp, out)
		}

		outRev := ciph.Decode(out)
		if !reflect.DeepEqual(outRev, c.in) {
			t.Fatalf("cipher reversal wrong. exp: %v, got: %v", c.in, outRev)
		}
	}
}

// should continue from the position of the whole stream, not the array
func TestCipherStream(t *testing.T) {
	ciphb := []byte{0x05}
	in1 := []byte{0x68, 0x65, 0x6c}
	exp1 := []byte{0x68, 0x66, 0x6e}
	in2 := []byte{0x6c, 0x6f}
	exp2 := []byte{0x6f, 0x73}

	ciph := cipher.ParseCipher(ciphb)
	out1 := ciph.Encode(in1)
	if !reflect.DeepEqual(out1, exp1) {
		t.Fatalf("cipher wrong. exp: %v, got: %v", exp1, out1)
	}
	out2 := ciph.Encode(in2)
	if !reflect.DeepEqual(out2, exp2) {
		t.Fatalf("cipher wrong. exp: %v, got: %v", exp2, out2)
	}
}
