package cipher

import (
	"bytes"
	"fmt"
	"io"
)

type cipherReader struct {
	r  io.Reader
	fn func([]byte) []byte
}

func (cr *cipherReader) Read(p []byte) (n int, err error) {
	n, err = cr.r.Read(p)
	if err != nil {
		return n, err
	}
	enc := cr.fn(p)
	copy(p, enc)
	return n, err
}

func ApplyCipherDecode(c Cipher, r io.Reader) io.Reader {
	return &cipherReader{
		r:  r,
		fn: c.Decode,
	}
}

type Cipher interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
}

type combinedCipher struct {
	operations []Cipher
}

func (cc *combinedCipher) Encode(b []byte) []byte {
	input := b
	for i, c := range cc.operations {
		input = c.Encode(input)
		fmt.Println(i, input)
	}
	return input
}

func (cc *combinedCipher) Decode(b []byte) []byte {
	input := b
	for i := len(cc.operations) - 1; i >= 0; i-- {
		c := cc.operations[i]
		input = c.Decode(input)
	}
	return input
}

func (cc *combinedCipher) String() string {
	return fmt.Sprint(cc.operations)
}

func ParseCipher(bs []byte) (Cipher, error) {
	input := bs
	operations := make([]Cipher, 0)
	copyForTesting := make([]Cipher, 0) // to prevent mutations

	for input != nil {
		var currByte byte
		currByte, input = getOne(input)

		switch currByte {
		case 0x01:
			operations = append(operations, reverseBit{})
			copyForTesting = append(copyForTesting, reverseBit{})
		case 0x02:
			var xorByte byte
			xorByte, input = getOne(input)
			operations = append(operations, xor{n: xorByte})
			copyForTesting = append(copyForTesting, xor{n: xorByte})
		case 0x03:
			operations = append(operations, &xorpos{})
			copyForTesting = append(copyForTesting, &xorpos{})
		case 0x04:
			var addByte byte
			addByte, input = getOne(input)
			operations = append(operations, add{n: addByte})
			copyForTesting = append(copyForTesting, add{n: addByte})
		case 0x05:
			operations = append(operations, &addpos{})
			copyForTesting = append(copyForTesting, &addpos{})
		}
	}

	testCipher := &combinedCipher{
		operations: copyForTesting,
	}
	if isCipherNoop(testCipher) {
		return nil, fmt.Errorf("cipher is noop")
	}

	return &combinedCipher{
		operations: operations,
	}, nil
}

// Test if cipher is noop.
// THIS MUTATES THE CIPHER
func isCipherNoop(c Cipher) bool {
	in := []byte{0x01, 0x02, 0x03, 0x04}
	out := c.Encode(in)
	return bytes.Equal(in, out)
}

func getOne[T any](arr []T) (T, []T) {
	first := arr[0]
	if len(arr) > 1 {
		return first, arr[1:]
	} else {
		return first, nil
	}
}

type reverseBit struct{}

func (rb reverseBit) Encode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		var num byte
		for i := 0; i < 8; i++ {
			val := (1 << i) & el
			if val != 0 {
				num |= (1 << (7 - i))
			}
		}
		ret[i] = num
	}
	return ret
}

func (rb reverseBit) Decode(bs []byte) []byte {
	return rb.Encode(bs) // inverse is itself :)
}

func (rb reverseBit) String() string {
	return "reversebits"
}

type xor struct {
	n byte
}

func (xr xor) Encode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		ret[i] = el ^ xr.n
	}
	return ret
}

func (xr xor) Decode(bs []byte) []byte {
	return xr.Encode(bs) // xor's inverse is itself :)
}

func (xr xor) String() string {
	return fmt.Sprintf("xor(%v)", xr.n)
}

type xorpos struct {
	encodePos int
	decodePos int
}

func (xrp *xorpos) Encode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		ret[i] = el ^ byte(xrp.encodePos)
		xrp.encodePos += 1
	}
	return ret
}

// xor's inverse is itself :)
func (xrp *xorpos) Decode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		ret[i] = el ^ byte(xrp.decodePos)
		xrp.decodePos += 1
	}
	return ret
}

func (xrp *xorpos) String() string {
	return fmt.Sprintf("xorpos(dec:%v,enc:%v)", xrp.decodePos, xrp.encodePos)
}

type add struct {
	n byte
}

func (ad add) Encode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) + int(ad.n)
		modded := added % 256
		ret[i] = byte(modded)
	}
	return ret
}

func (ad add) Decode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) - int(ad.n)
		modded := added % 256
		ret[i] = byte(modded)
	}
	return ret
}

func (ad add) String() string {
	return fmt.Sprintf("add(%v)", ad.n)
}

type addpos struct {
	encodePos int
	decodePos int
}

func (adp *addpos) Encode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) + adp.encodePos
		modded := added % 256
		ret[i] = byte(modded)
		adp.encodePos += 1
	}
	return ret
}

func (adp *addpos) Decode(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) - adp.decodePos
		modded := added % 256
		ret[i] = byte(modded)
		adp.decodePos += 1
	}
	return ret
}

func (adp *addpos) String() string {
	return fmt.Sprintf("addpos(dec:%v,enc:%v)", adp.decodePos, adp.encodePos)
}
