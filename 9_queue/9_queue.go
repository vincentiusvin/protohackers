package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"protohackers/9_queue/queue"
	"time"
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

	// instead of using json.NewDecoder, we chunk each request by newline.
	// easier error handling this way
	r := bufio.NewReader(c)

	enc := json.NewEncoder(c)

	for {
		b, err := r.ReadBytes('\n')
		t := time.Now()
		if err != nil {
			log.Println(err)
			return
		}

		val, err := queue.Decode(b)
		if err != nil {
			if err == io.EOF {
				return
			}

			log.Printf("decoding error: %v. handling gracefully...\n", err)

			err = enc.Encode(&queue.GeneralError{
				Status: queue.StatusError,
			})
			if err != nil {
				log.Println("failed to send error: %w", err)
				return
			}

			continue
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
		log.Printf("time: %v\n", time.Since(t))
	}
}
