package main

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	rw, in, out := rw()
	go handleIO(rw, "test")

	in <- "GET"
	rep := <-out
	if rep[:3] != "ERR" {
		t.Fatalf("expected error %v", rep)
	}
}

func rw() (rw *bufio.ReadWriter, in chan string, out chan string) {
	in = make(chan string, 1)
	out = make(chan string, 1)

	inr, inw := io.Pipe()
	r := bufio.NewReader(inr)

	go func() {
		defer inw.Close()
		for {
			d := <-in
			_, err := inw.Write([]byte(d + "\n"))
			if err != nil {
				return
			}
		}
	}()

	outr, outw := io.Pipe()
	w := bufio.NewWriter(outw)
	go func() {
		r := bufio.NewReader(outr)
		for {
			s, err := r.ReadString('\n')
			if err != nil {
				return
			}
			s = strings.TrimSpace(s)
			out <- s
		}
	}()

	rw = bufio.NewReadWriter(r, w)

	return
}
