package main

import "testing"

func TestBogus(t *testing.T) {
	type bogusCases struct {
		in       string
		expected string
	}
	cases := []bogusCases{
		{
			in:       "send money here: 7F1u3wSD5RbOHQmupo9nx4TnhQ",
			expected: "send money here: 7YWHMfk9JZe0LM0g1ZauHuiSxhI",
		},
	}

	for _, c := range cases {
		out := boguscoined(c.in)
		if out != c.expected {
			t.Fatalf("failed to transform boguscoin %v. expected %v. got %v", c.in, c.expected, out)
		}
	}

}
