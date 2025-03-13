package main

import (
	"bufio"
	"log"
	"net"
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
		log.Println(c.RemoteAddr(), "connected")
		go handleConnection(c)
	}
}

func connToChan(client net.Conn) chan string {
	cli_to_up := make(chan string)
	sc := bufio.NewScanner(client)
	go func() {
		for sc.Scan() {
			t := sc.Text()
			cli_to_up <- t
		}
		close(cli_to_up)
	}()
	return cli_to_up
}

func boguscoined(s string) string {
	return s
}

func handleConnection(client net.Conn) {
	defer client.Close()

	upstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		panic(err)
	}
	defer upstream.Close()

	cli_to_up := connToChan(client)
	up_to_cli := connToChan(upstream)

	for {
		select {
		case cli_to_up_msg, ok := <-cli_to_up:
			if !ok {
				break
			}
			doctored_msg := boguscoined(cli_to_up_msg)
			upstream.Write([]byte(doctored_msg))
		case up_to_cli_msg, ok := <-up_to_cli:
			if !ok {
				break
			}
			doctored_msg := boguscoined(up_to_cli_msg)
			client.Write([]byte(doctored_msg))
		}
	}
}
