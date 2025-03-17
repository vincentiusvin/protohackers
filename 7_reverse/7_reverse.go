package main

import (
	"log"
	"net"
	"strconv"
)

func main() {
	port := 8000
	c, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at port", strconv.Itoa(port))

	defer c.Close()

	for {
		ListenLRCP(c)
	}
}

func ListenLRCP(c *net.UDPConn) {
	b := make([]byte, 1000)
	_, addr, err := c.ReadFromUDP(b)
	if err != nil {
		panic(err)
	}
	log.Println("recv from: ", addr)
}
