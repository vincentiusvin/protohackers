package infra

import (
	"context"
	"encoding/binary"
	"io"
)

func EncodeMessages(ctx context.Context, w io.Writer) chan any {
	ch := make(chan any)

	go func() {
		defer close(ch)
		for msg := range ch {
			select {
			case <-ctx.Done():
				return
			default:
			}

			switch v := msg.(type) {
			case *Ticket:
				w.Write(EncodeTicket(v))
			case Heartbeat:
				w.Write(EncodeHeartbeat())
			case *SpeedError:
				w.Write(EncodeError(v))
			}
		}
	}()

	return ch
}

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

func EncodeError(err *SpeedError) []byte {
	return append([]byte{0x10}, EncodeString(err.Msg)...)
}

func EncodeString(s string) []byte {
	l := byte(uint8(len(s)))
	return append([]byte{l}, []byte(s)...)
}
