package lrcp_test

import (
	"protohackers/7_reverse/lrcp"
	"reflect"
	"testing"
)

func TestConnect(t *testing.T) {
	in := "/connect/1234567"
	exp := &lrcp.Connect{
		Session: 1234567,
	}
	c, err := lrcp.ParseConnect(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, exp) {
		t.Fatalf("expected %v got %v", exp, c.Session)
	}
}

func TestData(t *testing.T) {
	in := "/data/1234567/0/hello"
	exp := &lrcp.Data{
		Session: 1234567,
		Pos:     0,
		Data:    "hello",
	}
	c, err := lrcp.ParseData(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, exp) {
		t.Fatalf("expected %v got %v", exp, c.Session)
	}
}
