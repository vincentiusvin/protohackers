package git

import (
	"log"
	"slices"
	"testing"
)

func TestDirectory(t *testing.T) {
	dir := newDirectory("dir1")
	n := newFile("file1", []byte{0x01, 0x02})
	dir.addNode(n)
	log.Println(n)
}

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
