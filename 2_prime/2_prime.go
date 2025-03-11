package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		prime(c)
	}
}

type PrimeRequest struct {
	Method string `json:"method"`
	Number int    `json:"number"`
}

func prime(c net.Conn) {
	defer c.Close()

	sc := bufio.NewScanner(c)
	for sc.Scan() {
		s := sc.Text()
		handleRequest(s)
	}
	fmt.Println("EOF")
}

func handleRequest(s string) {
	b := new(bytes.Buffer)
	b.WriteString(s)

	d := json.NewDecoder(b)

	var pr PrimeRequest
	err := d.Decode(&pr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(pr)
}
