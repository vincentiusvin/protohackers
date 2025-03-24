package main

import (
	"bufio"
	"fmt"
	"io"
	"protohackers/10_git/git"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	_, in, out := rw()

	in <- "GET"
	rep := <-out
	if rep[:3] != "ERR" {
		t.Fatalf("expected error %v", rep)
	}

	in <- "GET /dir1 r1"
	fmt.Println(<-out)
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
	go handleIO(rw, "test", vc)

	return
}
