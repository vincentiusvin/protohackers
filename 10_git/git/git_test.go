package git

import (
	"slices"
	"testing"
)

func TestFileName(t *testing.T) {
	type fileNameCases struct {
		in  string
		exp []string
	}
	cases := []fileNameCases{
		{
			in:  "kucing/meong",
			exp: nil,
		},
		{
			in:  "/kucing//meong",
			exp: nil,
		},
		{
			in:  "/",
			exp: nil,
		},
		{
			in:  "/a",
			exp: []string{"a"},
		},
		{
			in:  "/kucing/meong",
			exp: []string{"kucing", "meong"},
		},
	}

	for _, c := range cases {
		out, _ := splitPaths(c.in)
		if !slices.Equal(out, c.exp) {
			t.Fatalf("failed to parse filename. exp %v got %v", c.exp, out)
		}
	}
}
