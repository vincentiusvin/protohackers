package main

import (
	"log"
	"net"
	"protohackers/6_speed/infra"
	"time"
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
		go handleConnection(c)
	}
}

type ClientType int

const (
	None ClientType = iota
	Camera
	Dispatcher
)

func handleConnection(c net.Conn) {
	defer c.Close()

	ch := infra.ParseMessages(c)
	send := func(b []byte) error {
		_, err := c.Write(b)
		return err
	}
	handleConnectionLogic(ch, send)
}

// With the connection abstracted away
func handleConnectionLogic(c chan any, send func([]byte) error) {
	clientType := None
	spawnedHb := false

	for msg := range c {
		switch v := msg.(type) {
		case *infra.IAmACamera:
			if clientType != None {
				send(infra.EncodeError("already another type"))
				continue
			}

			clientType = Camera
			log.Println("new camera", v)
		case *infra.IAmADispatcher:
			if clientType != None {
				send(infra.EncodeError("already another type"))
				continue
			}

			clientType = Dispatcher
			log.Println("new dispatcher", v)
		case *infra.WantHeartbeat:
			if spawnedHb {
				send(infra.EncodeError("already sending heartbeat"))
				continue
			}
			go func() {
				log.Println("starting heartbeat every", v.Interval, "ds")
				for {
					err := send(infra.EncodeHeartbeat())
					if err != nil {
						break
					}
					deciseconds := time.Second / 10
					time.Sleep(time.Duration(v.Interval) * deciseconds)
				}
				log.Println("stopped heartbeat every", v.Interval, "ds")
			}()
			spawnedHb = true
		case *infra.Plate:
			if clientType != Camera {
				send(infra.EncodeError("you are not a camera"))
				continue
			}
		}
	}
}
