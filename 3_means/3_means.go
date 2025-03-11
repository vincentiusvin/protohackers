package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go means(c)
	}
}

type Insert struct {
	Timestamp int
	Price     int
}

type Query struct {
	Mintime int
	Maxtime int
}

func means(c net.Conn) {
	defer c.Close()

	sess := make([]*Insert, 0)
	for {
		b := make([]byte, 9)
		_, err := io.ReadFull(c, b)
		if err != nil {
			log.Println(err)
			break
		}

		i, q := parsePacket(b)

		if i == nil && q == nil {
			log.Println("unable to parse", b)
		}

		if i != nil {
			sess = append(sess, i)
			log.Println("insert", i)
		}

		if q != nil {
			tally := 0
			count := 0
			for _, s := range sess {
				inside := (q.Mintime <= s.Timestamp) && (s.Timestamp <= q.Maxtime)
				if !inside {
					continue
				}

				tally += s.Price
				count += 1
			}

			var result uint32
			if count != 0 {
				result = uint32(tally / count)
			}
			resp := make([]byte, 4)
			binary.BigEndian.PutUint32(resp, result)
			c.Write(resp)
			log.Println("query", q, "resp:", result)
		}
	}
}

func parsePacket(b []byte) (*Insert, *Query) {
	t := b[0]

	// uint32 -> int32 -> int
	// casting directly to int will break negative numbers on 64 bit systems
	firstint := b[1:5]
	firstnum := int(int32(binary.BigEndian.Uint32(firstint)))

	secint := b[5:9]
	secnum := int(int32(binary.BigEndian.Uint32(secint)))

	if t == 'Q' {
		return nil, &Query{
			Mintime: firstnum,
			Maxtime: secnum,
		}
	} else if t == 'I' {
		return &Insert{
			Timestamp: firstnum,
			Price:     secnum,
		}, nil
	}
	return nil, nil
}
