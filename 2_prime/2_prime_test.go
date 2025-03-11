package main

import (
	"fmt"
	"testing"
)

type IsPrimeCases struct {
	in  float64
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
			in:  4,
			exp: false,
		},
		{
			in:  7,
			exp: true,
		},
		{
			in:  25,
			exp: false,
		},
		{
			in:  100,
			exp: false,
		},
		{
			in:  101,
			exp: true,
		},
		{
			in:  91028332887393654427978476145102147254121969766630449978,
			exp: false,
		},
	}

	for _, c := range cases {
		in := c.in
		exp := c.exp

		t.Run(fmt.Sprintf("%v", in), func(t *testing.T) {
			out := isPrime(in)

			if out != exp {
				t.Fatalf("wrong isprime for %v. exp %v got %v", in, exp, out)
			}
		})

	}

}
