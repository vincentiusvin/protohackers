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
	}

	for _, c := range cases {
		out := cipher.RunCipher(c.ciph, c.in)
		if !reflect.DeepEqual(out, c.exp) {
			t.Fatalf("cipher wrong. exp: %v, got: %v", c.exp, out)
		}
	}
}
