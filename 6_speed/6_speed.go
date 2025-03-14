package main

import (
	"context"
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

	in := infra.ParseMessages(c)
	out := infra.EncodeMessages(nil, c)
	handleConnectionLogic(nil, in, out)
}

// With the connection abstracted away
func handleConnectionLogic(ctx context.Context, incoming chan any, outgoing chan any) {
	clientType := None
	spawnedHb := false

	sendError := func(msg string) {
		outgoing <- infra.SpeedError{
			Msg: msg,
		}
	}

	for msg := range incoming {
		switch v := msg.(type) {
		case *infra.IAmACamera:
			if clientType != None {
				sendError("already another type")
				continue
			}

			clientType = Camera
			log.Println("new camera", v)
		case *infra.IAmADispatcher:
			if clientType != None {
				sendError("already another type")
				continue
			}

			clientType = Dispatcher
			log.Println("new dispatcher", v)
		case *infra.WantHeartbeat:
			if spawnedHb {
				sendError("already sending heartbeat")
				continue
			}

			go func() {
				log.Println("starting heartbeat every", v.Interval, "ds")
				for {
					select {
					case <-ctx.Done():
						log.Println("stopped heartbeat every", v.Interval, "ds")
						return
					case outgoing <- infra.Heartbeat{}:
						deciseconds := time.Second / 10
						time.Sleep(time.Duration(v.Interval) * deciseconds)
					}
				}
			}()
			spawnedHb = true
		case *infra.Plate:
			if clientType != Camera {
				sendError("you are not a camera")
				continue
			}
		}
	}
}
