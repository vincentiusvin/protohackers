package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
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

	clientNum := jc.GetClientID()

	dr := &queue.DisconnectRequest{
		ClientID: clientNum,
	}
	defer jc.DisconnectWorker(dr)

	dec := json.NewDecoder(c)
	enc := json.NewEncoder(c)

	for {
		val, err := queue.Decode(dec)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Println(err)
				err = enc.Encode(&queue.GeneralError{
					Status: queue.StatusError,
				})
				if err != nil {
					log.Println("failed to send error: %w", err)
				}
				continue
			}
		}

		switch t := val.(type) {
		case *queue.GetRequest:
			t.ClientID = clientNum
			resp := jc.Get(t)
			err = enc.Encode(resp)
		case *queue.PutRequest:
			resp := jc.Put(t)
			err = enc.Encode(resp)
		case *queue.DeleteRequest:
			resp := jc.Delete(t)
			err = enc.Encode(resp)
		case *queue.AbortRequest:
			t.ClientID = clientNum
			resp := jc.Abort(t)
			err = enc.Encode(resp)
		}

		if err != nil {
			log.Println(err)
		}
	}
}
