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

	pc := pest.NewControllerTCP()

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
	log.Printf("[%v] got hello\n", addr)

	helloB := infra.Encode(types.Hello{
		Protocol: "pestcontrol",
		Version:  1,
	})

	_, err := c.Write(helloB)
	if err != nil {
		return
	}
	log.Printf("[%v] sent hello\n", addr)

	for v := range visitChan {
		log.Printf("[%v] added visit: %v\n", addr, v)
		err = pc.AddSiteVisit(v)
		if err != nil {
			log.Printf("%v got err %v", addr, err)
			break
		}
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
			n, err := c.Read(b)
			curr = append(curr, b[:n]...)
			if err != nil {
				break
			}

			for {
				res := infra.Parse(curr)
				if !res.Ok {
					if len(curr) != 0 {
						log.Printf("HELP: %v\n", curr)
					}
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
