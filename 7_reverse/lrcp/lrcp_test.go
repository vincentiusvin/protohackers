package lrcp_test

import (
	"protohackers/7_reverse/lrcp"
	"testing"
)

func TestConnect(t *testing.T) {
	in := "/connect/1234567"
	var exp uint = 1234567
	c, err := lrcp.ParseConnect(in)
	if err != nil {
		t.Fatal(err)
	}
	if c.Session != exp {
		t.Fatalf("expected %v got %v", exp, c.Session)
	}
}
