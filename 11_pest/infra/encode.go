package infra

import (
	"encoding/binary"
	"protohackers/11_pest/pest"
)

func Encode(v any) []byte {
	switch nv := v.(type) {
	case pest.Hello:
		return encodeHello(nv)
	}
	return nil
}

func encodeHello(h pest.Hello) []byte {
	body := make([]byte, 0)
	body = append(body, encodeString(h.Protocol)...)
	body = append(body, encodeUint32(h.Version)...)

	return encaseEnvelope(body, 0x50)
}

func encaseEnvelope(b []byte, prefix byte) (ret []byte) {
	// total len = prefix (1 byte) + msgLen (4 byte) + msg (len(b) bytes) + checksum (1 byte)
	// 1 + 4 + len(b) + 1
	// len(b) + 6
	totalLen := uint32(len(b) + 6)

	ret = append(ret, prefix)
	ret = append(ret, encodeUint32(totalLen)...)
	ret = append(ret, b...)

	var sum byte = 0
	for i := 0; i < len(ret); i++ {
		sum += ret[i]
	}
	checkSum := 255 - sum + 1
	ret = append(ret, checkSum)

	return ret
}

func encodeUint32(val uint32) (ret []byte) {
	ret = binary.BigEndian.AppendUint32(ret, val)
	return ret
}

func encodeString(s string) (ret []byte) {
	strlen := uint32(len(s))
	ret = append(ret, encodeUint32(strlen)...)
	ret = append(ret, []byte(s)...)
	return ret
}
