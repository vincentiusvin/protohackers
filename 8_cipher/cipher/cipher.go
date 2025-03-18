package cipher

import "fmt"

type operation interface {
	execute([]byte) []byte
}

func RunCipher(ciph []byte, input []byte) []byte {
	fns := parseCipher(ciph)
	for _, c := range fns {
		input = c.execute(input)
	}
	return input
}

func parseCipher(bs []byte) []operation {
	input := bs
	operations := make([]operation, 0)

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
			operations = append(operations, xorpos{})
		case 0x04:
			var addByte byte
			addByte, input = getOne(input)
			operations = append(operations, add{n: addByte})
		case 0x05:
			operations = append(operations, addpos{})
		}
	}

	return operations
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

func (rb reverseBit) execute(bs []byte) []byte {
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

func (rb reverseBit) String() string {
	return "reversebits"
}

type xor struct {
	n byte
}

func (xr xor) execute(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		ret[i] = el ^ xr.n
	}
	return ret
}

func (xr xor) String() string {
	return fmt.Sprintf("xor(%v)", xr.n)
}

type xorpos struct{}

func (xrp xorpos) execute(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		ret[i] = el ^ byte(i)
	}
	return ret
}

func (xrp xorpos) String() string {
	return "xorpos"
}

type add struct {
	n byte
}

func (ad add) execute(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) + int(ad.n)
		modded := added % 256
		ret[i] = byte(modded)
	}
	return ret
}

func (ad add) String() string {
	return fmt.Sprintf("add(%v)", ad.n)
}

type addpos struct{}

func (adp addpos) execute(bs []byte) []byte {
	ret := make([]byte, len(bs))
	for i, el := range bs {
		added := int(el) + i
		modded := added % 256
		ret[i] = byte(modded)
	}
	return ret
}

func (adp addpos) String() string {
	return "addpos"
}
