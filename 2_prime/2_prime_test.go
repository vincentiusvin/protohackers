package main

import (
	"strconv"
	"testing"
)

type IsPrimeCases struct {
	in  int
	exp bool
}

func TestIsPrime(t *testing.T) {
	cases := []IsPrimeCases{
		{
			in:  -4,
			exp: false,
		},
		{
			in:  -2,
			exp: false,
		},
		{
			in:  1,
			exp: false,
		},
		{
			in:  2,
			exp: true,
		},
		{
			in:  3,
			exp: true,
		},
		{
			in:  100,
			exp: false,
		},
		{
			in:  101,
			exp: true,
		},
	}

	for _, c := range cases {
		in := c.in
		exp := c.exp

		t.Run(strconv.Itoa(in), func(t *testing.T) {
			out := isPrime(in)

			if out != exp {
				t.Fatalf("wrong isprime for %v. exp %v got %v", in, exp, out)
			}
		})

	}

}
