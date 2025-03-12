package main

import (
	"testing"
)

func TestNameValidate(t *testing.T) {
	type nameCases struct {
		name string
		ok   bool
	}

	cases := []nameCases{
		{
			name: "udin",
			ok:   true,
		},
		{
			name: "udin_swag",
			ok:   false,
		},
		{
			name: "UdinSwag100",
			ok:   true,
		},
	}

	for _, c := range cases {
		out := validateName(c.name)
		if out != c.ok {
			t.Fatalf("wrong output for %v. expect: %v got: %v", c.name, c.ok, out)
		}
	}
}
