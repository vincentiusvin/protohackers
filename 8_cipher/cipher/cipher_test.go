package cipher_test

import (
	"protohackers/8_cipher/cipher"
	"testing"
)

func TestCipher(t *testing.T) {
	in := []byte{0x02, 0x00, 0x00}
	cipher.ParseCipher(in)
}
