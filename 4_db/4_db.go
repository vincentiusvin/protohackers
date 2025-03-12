package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

// You need to reply on the same port as you are listening.
// Tried using Dial, didn't work

func main() {
	addr := ":8000"
	ln, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	m := make(map[string]string)
	m["version"] = "database punya udin 1.0"

	for {
		b := make([]byte, 1000)
		_, addr, err := ln.ReadFrom(b)
		if err != nil {
			panic(err)
		}
		log.Println("recv from: ", addr)

		// "udin\0" != "udin"
		// gotta remove the nulls first
		b = bytes.Trim(b, "\000")

		ins, ret := ParseRequest(string(b))

		if ins != nil && ins.Key != "version" {
			m[ins.Key] = ins.Value
			log.Println("ins for", ins.Key, ":", ins.Value)
		}

		if ret != nil {
			retval := fmt.Sprintf("%v=%v", ret.Key, m[ret.Key])

			c, err := net.Dial("udp", addr.String())
			if err != nil {
				panic(err)
			}
			c.Write([]byte(retval))
			c.Close()

			log.Println("ret for", ret.Key, ":", m[ret.Key])
		}
	}
}

type Insert struct {
	Key   string
	Value string
}

type Retrieve struct {
	Key string
}

func ParseRequest(s string) (*Insert, *Retrieve) {
	bef, aft, found := strings.Cut(s, "=")

	if found {
		ins := new(Insert)
		ins.Key = bef
		ins.Value = aft
		return ins, nil
	}

	ret := new(Retrieve)
	ret.Key = s

	return nil, ret
}
