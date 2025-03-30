package main

import (
	"log"
	"net"
	"protohackers/11_pest/infra"
	"protohackers/11_pest/pest"
	"protohackers/11_pest/types"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Server listening at " + addr)

	defer ln.Close()

	pc := pest.NewController()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(c, pc)
	}
}

func handleConn(c net.Conn, pc pest.Controller) {
	defer c.Close()

	addr := c.RemoteAddr().String()
	log.Printf("[%v] connected\n", addr)

	helloChan, visitChan := processIncoming(c)

	helloReply := <-helloChan
	if helloReply.Protocol != "pestcontrol" || helloReply.Version != 1 {
		return
	}

	for v := range visitChan {
		log.Printf("[%v] added visit\n", v)
		pc.AddSiteVisit(v)
	}
}

func processIncoming(c net.Conn) (helloChan chan types.Hello, visitChan chan types.SiteVisit) {
	helloChan = make(chan types.Hello)
	visitChan = make(chan types.SiteVisit)

	go func() {
		defer func() {
			close(helloChan)
			close(visitChan)
		}()

		var curr []byte
		for {
			b := make([]byte, 1024)
			_, err := c.Read(b)
			curr = append(curr, b...)
			if err != nil {
				break
			}

			for {
				res := infra.Parse(curr)
				if !res.Ok {
					break
				}

				switch v := res.Value.(type) {
				case types.Hello:
					helloChan <- v
				case types.SiteVisit:
					visitChan <- v
				}

				curr = res.Next
			}
		}
	}()

	return helloChan, visitChan
}
