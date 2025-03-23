package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"protohackers/9_queue/queue"
)

func main() {
	addr := ":8000"
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	jc := queue.NewJobCenter(ctx)

	log.Println("Server listening at " + addr)

	defer ln.Close()

	for {
		c, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c, jc)
	}
}

func handleConnection(c net.Conn, jc *queue.JobCenter) {
	defer c.Close()

	clientNum := rand.Int()

	dr := &queue.DisconnectRequest{
		ClientID: clientNum,
	}
	defer jc.DisconnectWorker(dr)

}
