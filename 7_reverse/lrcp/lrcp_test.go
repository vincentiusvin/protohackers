package lrcp_test

import (
	"protohackers/7_reverse/lrcp"
	"reflect"
	"testing"
)

func TestParsing(t *testing.T) {
	type parseCases struct {
		in  string
		exp lrcp.LRCPPackets
	}

	runParse := func(t *testing.T, cs []parseCases) {
		for _, c := range cs {
			out, _ := lrcp.Parse(c.in)
			if !reflect.DeepEqual(out, c.exp) {
				t.Fatalf("expected %v got %v", c.exp, out)
			}
			if c.exp == nil {
				continue
			}
			reencode := c.exp.Encode()
			if reencode != c.in {
				t.Fatalf("expected %v got %v", c.in, reencode)
			}
		}
	}

	t.Run("connect cases", func(t *testing.T) {
		connectCases := []parseCases{
			{
				in: "/connect/1234567/",
				exp: &lrcp.Connect{
					Session: 1234567,
				},
			},
			{
				in:  "/connect/1234567/ayambakar",
				exp: nil,
			},
		}
		runParse(t, connectCases)
	})

	t.Run("ack cases", func(t *testing.T) {
		ackCases := []parseCases{
			{
				in: "/ack/1234567/1024/",
				exp: &lrcp.Ack{
					Session: 1234567,
					Length:  1024,
				},
			},
			{
				in:  "/ack/1234567/nasigoreng",
				exp: nil,
			},
		}
		runParse(t, ackCases)
	})

	t.Run("data cases", func(t *testing.T) {
		dataCases := []parseCases{
			{
				in: "/data/1234567/0//",
				exp: &lrcp.Data{
					Session: 1234567,
					Pos:     0,
					Data:    "",
				},
			},
			{
				in: "/data/1234567/0/hello/",
				exp: &lrcp.Data{
					Session: 1234567,
					Pos:     0,
					Data:    "hello",
				},
			},
			{
				in:  "/data/1234567/0/hello",
				exp: nil,
			},
			{
				in:  "/data/1234567/0/hello/darkness/my/",
				exp: nil,
			},
			{
				in: "/data/1234567/0/foo\\/bar\\\\baz/",
				exp: &lrcp.Data{
					Session: 1234567,
					Pos:     0,
					Data:    "foo/bar\\baz",
				},
			},
		}
		runParse(t, dataCases)
	})

	t.Run("close cases", func(t *testing.T) {
		closeCases := []parseCases{
			{
				in: "/close/1234567/",
				exp: &lrcp.Close{
					Session: 1234567,
				},
			},
			{
				in:  "/close/1234567",
				exp: nil,
			},
		}
		runParse(t, closeCases)
	})
}

type lrcpMock struct {
	in  chan string
	out chan string
}

func (lm *lrcpMock) Read() ([]byte, error) {
	b := <-lm.in
	return []byte(b), nil
}

func (lm *lrcpMock) Write(b []byte) error {
	lm.out <- string(b)
	return nil
}

func TestLRCPReceive(t *testing.T) {
	ls := lrcp.MakeLRCPServer()

	chin := make(chan string)
	chout := make(chan string)

	ls.Listen(func() lrcp.LRCPListenerSession {
		return &lrcpMock{
			in:  chin,
			out: chout,
		}
	})

	final := make(chan string)

	go func() {
		sess := ls.Accept()
		acc := ""
		outch, err := sess.Resolve()
		if err != nil {
			panic("lrcp session channel errored unexpectedly")
		}
		for s := range outch {
			acc += s
		}
		final <- acc
	}()

	chin <- "/connect/1234567/"
	if <-chout != "/ack/1234567/0/" {
		t.Fatalf("wrong reply to connect")
	}
	chin <- "/data/1234567/0/hello/"
	if <-chout != "/ack/1234567/5/" {
		t.Fatalf("wrong reply to data 1")
	}
	chin <- "/data/1234567/5/meong/"
	if <-chout != "/ack/1234567/10/" {
		t.Fatalf("wrong reply to data 2")
	}
	chin <- "/data/1234567/10/123/"
	if <-chout != "/ack/1234567/13/" {
		t.Fatalf("wrong reply to data 3")
	}
	chin <- "/close/1234567/"
	if <-chout != "/close/1234567/" {
		t.Fatalf("wrong reply to close")
	}
	if <-final != "hellomeong123" {
		t.Fatalf("wrong assembled string")
	}
}

func TestLRCPSend(t *testing.T) {
	ls := lrcp.MakeLRCPServer()

	chin := make(chan string)
	chout := make(chan string)

	ls.Listen(func() lrcp.LRCPListenerSession {
		return &lrcpMock{
			in:  chin,
			out: chout,
		}
	})

	go func() {
		sess := ls.Accept()
		sess.SendData("meong")
	}()

	chin <- "/connect/1234567/"
	if <-chout != "/ack/1234567/0/" {
		t.Fatalf("wrong reply to connect")
	}

	if <-chout != "/data/1234567/0/meong/" {
		t.Fatalf("wrong reply to connect")
	}

	chin <- "/close/1234567/"
	if <-chout != "/close/1234567/" {
		t.Fatalf("wrong reply to close")
	}
}
