package cipher

import "fmt"

type Cipher interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
}

type combinedCipher struct {
	operations []Cipher
}

func (cc *combinedCipher) Encode(b []byte) []byte {
	input := b
	for _, c := range cc.operations {
		input = c.Encode(input)
	}
	return input
}

func (cc *combinedCipher) Decode(b []byte) []byte {
	input := b
	for _, c := range cc.operations {
		input = c.Decode(input)
	}
	return input
}

func ParseCipher(bs []byte) Cipher {
	input := bs
	operations := make([]Cipher, 0)

	for input != nil {
		var currByte byte
		currByte, input = getOne(input)

		switch currByte {
		case 0x01:
			operations = append(operations, reverseBit{})
		case 0x02:
			var xorByte byte
			xorByte, input = getOne(input)
			operations = append(operations, xor{n: xorByte})
		case 0x03:
			operations = append(operations, &xorpos{})
		case 0x04:
			var addByte byte
			addByte, input = getOne(input)
			operations = append(operations, add{n: addByte})
		case 0x05:
			operations = append(operations, &addpos{})
		}
	}

	return &combinedCipher{
		operations: operations,
	}
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
	return "xorpos"
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
	return "addpos"
}
