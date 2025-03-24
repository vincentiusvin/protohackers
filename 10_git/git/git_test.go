package git_test

import (
	"protohackers/10_git/git"
	"testing"
)

func TestPut(t *testing.T) {
	type putCases struct {
		in string
		ok bool
	}

	cases := []putCases{
		{
			in: "/dir1/dir2/file",
			ok: true,
		},
		{
			in: "/dir1/dir2/file",
			ok: true,
		},
		{
			in: "kucing/meong",
			ok: false,
		},
		{
			in: "/kucing//meong",
			ok: false,
		},
		{
			in: "/",
			ok: false,
		},
		{
			in: "/a",
			ok: true,
		},
		{
			in: "/kucing/meong",
			ok: true,
		},
	}

	for _, c := range cases {
		t.Run("put", func(t *testing.T) {
			v := git.NewVersionControl()
			_, err := v.PutFile(c.in, []byte{0x01})
			if c.ok {
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if err == nil {
					t.Fatalf("expected error")
				}
			}
		})
	}
}
