package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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

func connToChan(conn net.Conn) chan string {
	cli_to_up := make(chan string)
	sc := bufio.NewScanner(conn)
	go func() {
		for sc.Scan() {
			t := sc.Text()
			cli_to_up <- t
		}
		fmt.Println("closed")
		close(cli_to_up)
	}()
	return cli_to_up
}

func isAlphanum(s string) bool {
	for _, c := range s {
		uppercase := (c >= 'A') && (c <= 'Z')
		lowercase := (c >= 'a') && (c <= 'z')
		digits := (c >= '0') && (c <= '9')

		if uppercase || lowercase || digits {
			continue
		}

		return false
	}
	return true
}

func boguscoined(s string) string {
	tonycoin := "7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	splits := strings.Split(s, " ")
	rets := make([]string, 0)
	for _, split := range splits {
		if len(split) < 26 || len(split) > 35 {
			rets = append(rets, split)
			continue
		}
		if split[0] != '7' {
			rets = append(rets, split)
			continue
		}

		if !isAlphanum(split) {
			rets = append(rets, split)
			continue
		}

		rets = append(rets, tonycoin)
	}
	return strings.Join(rets, " ")
}

func handleConnection(client net.Conn) {
	log.Println(client.RemoteAddr(), "connected")
	defer func() {
		client.Close()
		log.Println(client.RemoteAddr(), "disconnected")
	}()

	upstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		panic(err)
	}
	log.Println("connected to upstream")
	defer func() {
		log.Println("upstream disconnected")
		upstream.Close()
	}()

	cli_to_up := connToChan(client)
	up_to_cli := connToChan(upstream)

	for {
		select {
		case cli_to_up_msg, ok := <-cli_to_up:
			if !ok {
				return
			}
			fmt.Println("client say:", cli_to_up_msg)
			doctored_msg := boguscoined(cli_to_up_msg) + "\n"
			_, err := upstream.Write([]byte(doctored_msg))
			if err != nil {
				panic(err)
			}
			fmt.Println("relayed to upstream:", doctored_msg)

		case up_to_cli_msg, ok := <-up_to_cli:
			if !ok {
				return
			}
			fmt.Println("upstream say:", up_to_cli_msg)
			doctored_msg := boguscoined(up_to_cli_msg) + "\n"
			_, err := client.Write([]byte(doctored_msg))
			if err != nil {
				panic(err)
			}
			fmt.Println("relayed to client:", doctored_msg)
		}
	}
}
