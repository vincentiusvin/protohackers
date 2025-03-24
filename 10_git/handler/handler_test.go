package handler_test

import (
	"bufio"
	"io"
	"protohackers/10_git/git"
	"protohackers/10_git/handler"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	_, in, out := rw()

	consumeReady := func() {
		r := <-out
		if r != "READY" {
			t.Fatalf("expected ready got %v", r)
		}
	}

	consumeReady()
	in <- "GET"
	rep := <-out
	if rep[:3] != "ERR" {
		t.Fatalf("expected error %v", rep)
	}

	consumeReady()
	in <- "PUT /dir1/file1 7\nkucing" // 7 since our writer appends a \n
	rep = <-out
	if rep != "OK r1" {
		t.Fatalf("expected success %v", rep)
	}

	consumeReady()
	in <- "GET /dir1/file1 r1"
	rep = <-out
	if rep != "OK 7" {
		t.Fatalf("expected success %v", rep)
	}
	data := <-out
	if data != "kucing" {
		t.Fatalf("expected kucing got %v", data)
	}

	consumeReady()
	in <- "LIST /dir1"
	rep = <-out
	if rep != "OK 1" {
		t.Fatalf("expected success %v", rep)
	}
	files := <-out
	if files != "file1 r1" {
		t.Fatalf("expected file1 r1 got %v", files)
	}
}

func rw() (vc *git.VersionControl, in chan string, out chan string) {
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

	rw := bufio.NewReadWriter(r, w)
	vc = git.NewVersionControl()
	go handler.HandleIO(rw, "test", vc)

	return
}
