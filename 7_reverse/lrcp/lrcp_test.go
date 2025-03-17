package lrcp_test

import (
	"protohackers/7_reverse/lrcp"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	in := "/connect/1234567/"
	exp := &lrcp.Connect{
		Session: 1234567,
	}
	c, err := lrcp.ParseConnect(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, exp) {
		t.Fatalf("expected %v got %v", exp, c)
	}
	reencode := exp.Encode()
	if reencode != in {
		t.Fatalf("expected %v got %v", in, reencode)
	}
}

func TestData(t *testing.T) {
	type dataCases struct {
		in  string
		exp *lrcp.Data
	}

	cases := []dataCases{
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
			in: "/data/1234567/0/foo\\/bar\\\\baz/",
			exp: &lrcp.Data{
				Session: 1234567,
				Pos:     0,
				Data:    "foo/bar\\baz",
			},
		},
	}

	for _, testCase := range cases {
		c, err := lrcp.ParseData(testCase.in)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(c, testCase.exp) {
			t.Fatalf("expected %v got %v", testCase.exp, c)
		}
		reencode := testCase.exp.Encode()
		if reencode != testCase.in {
			t.Fatalf("expected %v got %v", testCase.in, reencode)
		}
	}
}

func TestAck(t *testing.T) {
	in := "/ack/1234567/1024/"
	exp := &lrcp.Ack{
		Session: 1234567,
		Length:  1024,
	}
	c, err := lrcp.ParseAck(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, exp) {
		t.Fatalf("expected %v got %v", exp, c)
	}
	reencode := exp.Encode()
	if reencode != in {
		t.Fatalf("expected %v got %v", in, reencode)
	}
}

func TestClose(t *testing.T) {
	in := "/close/1234567/"
	exp := &lrcp.Close{
		Session: 1234567,
	}
	c, err := lrcp.ParseClose(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, exp) {
		t.Fatalf("expected %v got %v", exp, c)
	}
	reencode := exp.Encode()
	if reencode != in {
		t.Fatalf("expected %v got %v", in, reencode)
	}
}
