package infra

import (
	"context"
	"encoding/binary"
	"io"
)

type Encode interface {
	Encode() []byte
}

func EncodeMessages(ctx context.Context, w io.Writer) chan Encode {
	ch := make(chan Encode)

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-ch:
				w.Write(v.Encode())
			}

		}
	}()

	return ch
}

func (t *Ticket) Encode() []byte {
	ret := append(
		[]byte{0x21},
		encodeString(t.Plate)...,
	)
	ret = binary.BigEndian.AppendUint16(ret, t.Road)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile1)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp1)
	ret = binary.BigEndian.AppendUint16(ret, t.Mile2)
	ret = binary.BigEndian.AppendUint32(ret, t.Timestamp2)
	ret = binary.BigEndian.AppendUint16(ret, t.Speed)
	return ret
}

func (Heartbeat) Encode() []byte {
	return []byte{0x41}
}

func (err *SpeedError) Encode() []byte {
	return append([]byte{0x10}, encodeString(err.Msg)...)
}

func encodeString(s string) []byte {
	l := byte(uint8(len(s)))
	return append([]byte{l}, []byte(s)...)
}
