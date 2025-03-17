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
			_, addr, err := c.ReadFromUDP(b)
			if err != nil {
				panic(err)
			}
			log.Println("recv from: ", addr)

			// string(b[n])

		}
	}()

	return ret
}
