package infra

import "encoding/binary"

func EncodeTicket(t *Ticket) []byte {
	ret := append(
		[]byte{0x21},
		EncodeString(t.Plate)...,
	)
	ret = binary.BigEndian.AppendUint16(ret, t.Road)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile1)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp1)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile2)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp2)
	ret = binary.BigEndian.AppendUint16(ret, t.Speed)
	return ret
}

func EncodeHeartbeat() []byte {
	return []byte{0x41}
}

func EncodeError(err string) []byte {
	return append([]byte{0x10}, EncodeString(err)...)
}

func EncodeString(s string) []byte {
	l := byte(uint8(len(s)))
	return append([]byte{l}, []byte(s)...)
}
