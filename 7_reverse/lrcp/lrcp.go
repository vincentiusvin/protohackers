// provide a tcp like interface
// LRCP(in chan []byte) *LRCP
// LRCP.Accept() returns a session

package lrcp

import (
	"log"
	"net"
)

type LRCPServer struct {
}

func ParseUDP(c *net.UDPConn) chan any {
	ret := make(chan any)

	go func() {
		for {
			b := make([]byte, 1000)
			n, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			log.Println("recv from: ", addr)

			str := string(b[n])

			conn, err := ParseConnect(str)
			if err == nil {
				ret <- conn
			}

			data, err := ParseData(str)
			if err == nil {
				ret <- data
			}

			ack, err := ParseAck(str)
			if err == nil {
				ret <- ack
			}

			cls, err := ParseClose(str)
			if err == nil {
				ret <- cls
			}

			panic("failed to parse")
		}
	}()

	return ret
}
