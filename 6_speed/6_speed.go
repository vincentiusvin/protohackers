package main

import (
	"context"
	"log"
	"net"
	"protohackers/6_speed/infra"
	"protohackers/6_speed/ticketing"
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

	ctrl := ticketing.MakeController()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c, ctrl)
	}
}

type ClientType int

const (
	None ClientType = iota
	Camera
	Dispatcher
)

func handleConnection(c net.Conn, ctrl *ticketing.Controller) {
	defer c.Close()

	log.Println(c.RemoteAddr(), "connected")
	ctx, cancel := context.WithCancel(context.Background())

	in := infra.ParseMessages(c, cancel) // this guy cancels the other two
	out := infra.EncodeMessages(ctx, c)
	handleConnectionLogic(ctx, in, out, ctrl)
}

// With the connection abstracted away
func handleConnectionLogic(ctx context.Context, incoming chan any, outgoing chan infra.Encode, ctrl *ticketing.Controller) {
	clientType := None
	spawnedHb := false
	var camera *infra.IAmACamera

	sendError := func(msg string) {
		outgoing <- &infra.SpeedError{
			Msg: msg,
		}
	}

	for msg := range incoming {
		switch v := msg.(type) {
		case *infra.SpeedError:
			sendError("this is a server message")

		case *infra.Plate:
			if clientType != Camera || camera == nil {
				sendError("you are not a camera")
				continue
			}
			ctrl.AddPlates(&ticketing.Plate{
				Plate:     v.Plate,
				Timestamp: v.Timestamp,
				Road:      camera.Road,
				Mile:      camera.Mile,
			})

		case *infra.Ticket:
			sendError("this is a server message")

		case *infra.WantHeartbeat:
			if spawnedHb {
				sendError("already sending heartbeat")
				continue
			}

			if v.Interval == 0 {
				continue
			}

			go func() {
				log.Println("starting heartbeat every", v.Interval, "ds")
				for {
					select {
					case <-ctx.Done():
						log.Println("stopped heartbeat every", v.Interval, "ds")
						close(outgoing)
						return
					case outgoing <- infra.Heartbeat{}:
						deciseconds := time.Second / 10
						time.Sleep(time.Duration(v.Interval) * deciseconds)
					}
				}
			}()
			spawnedHb = true

		case infra.Heartbeat:
			sendError("this is a server message")

		case *infra.IAmACamera:
			if clientType != None {
				sendError("already another type")
				continue
			}
			ctrl.UpdateLimit(v.Road, v.Limit)
			clientType = Camera
			camera = v
			log.Printf("new camera on road %v mile %v\n", v.Road, v.Mile)

		case *infra.IAmADispatcher:
			if clientType != None {
				sendError("already another type")
				continue
			}

			ch := make(chan *infra.Ticket)
			go func() {
				for c := range ch {
					outgoing <- c
				}
			}()
			ctrl.AddDispatcher(v.Roads, ch)
			clientType = Dispatcher
			log.Printf("new dispatcher on roads %v\n", v.Roads)
		}
	}
}
