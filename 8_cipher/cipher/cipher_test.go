package cipher_test

import (
	"bufio"
	"bytes"
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
		{
			ciph: []byte{0x02, 0x01, 0x01, 0x00},
			in:   []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			exp:  []byte{0x96, 0x26, 0xb6, 0xb6, 0x76},
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
	inv1 := ciph.Decode(out1)
	if !reflect.DeepEqual(inv1, in1) {
		t.Fatalf("cipher wrong. exp: %v, got: %v", in1, inv1)
	}
	out2 := ciph.Encode(in2)
	if !reflect.DeepEqual(out2, exp2) {
		t.Fatalf("cipher wrong. exp: %v, got: %v", exp2, out2)
	}
	inv2 := ciph.Decode(out2)
	if !reflect.DeepEqual(inv2, in2) {
		t.Fatalf("cipher wrong. exp: %v, got: %v", in2, inv2)
	}
}

func TestReader(t *testing.T) {
	ciphB := []byte{0x02, 0x7b, 0x05, 0x01, 0x00}
	in := []byte{0xf2, 0x20, 0xba, 0x44, 0x18, 0x84, 0xba, 0xaa, 0xd0, 0x26, 0x44, 0xa4, 0xa8, 0x7e}
	exp := "4x dog,5x car\n"

	ciph := cipher.ParseCipher(ciphB)
	inBuf := bytes.NewBuffer(in)
	decodedIn := cipher.ApplyCipherDecode(ciph, inBuf)
	decodedR := bufio.NewReader(decodedIn)
	out, err := decodedR.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}

	if exp != out {
		t.Fatalf("cipher reader wrong result. exp %v got %v", exp, out)
	}
}

func TestReaderNoop(t *testing.T) {
	ciphB := []byte{0x00}
	in := []byte("cat\n")

	ciph := cipher.ParseCipher(ciphB)
	inBuf := bytes.NewBuffer(in)
	decodedIn := cipher.ApplyCipherDecode(ciph, inBuf)
	decodedR := bufio.NewReader(decodedIn)
	_, err := decodedR.ReadString('\n')
	if err == nil {
		t.Fatal("expected an error for noop ciphers")
	}
}
