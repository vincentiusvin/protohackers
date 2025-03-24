package git_test

import (
	"bytes"
	"protohackers/10_git/git"
	"testing"
)

type fixtureRet struct {
	v    *git.VersionControl
	f1   string
	f2   string
	f1b1 []byte
	f1b2 []byte
	f2b1 []byte
}

func vcFixture() fixtureRet {
	v := git.NewVersionControl()
	f1b1 := []byte{0x01, 0x02}
	f1b2 := []byte{0x01, 0x04}
	f2b1 := []byte{0x01, 0x03}
	f1 := "/dir1/dir2/file"
	f2 := "/dir1/dir2"
	v.PutFile(f1, f1b1)
	v.PutFile(f1, f1b2)
	v.PutFile(f2, f2b1)
	return fixtureRet{
		v:    v,
		f1b1: f1b1,
		f1b2: f1b2,
		f2b1: f2b1,
		f1:   f1,
		f2:   f2,
	}
}

func TestGet(t *testing.T) {
	ex := vcFixture()
	type getCases struct {
		inpath string
		inRev  int
		expB   []byte
		expErr bool
	}
	cases := []getCases{
		{
			inpath: ex.f1,
			inRev:  0,
			expErr: false,
			expB:   ex.f1b2,
		},
		{
			inpath: ex.f1,
			inRev:  2,
			expErr: false,
			expB:   ex.f1b2,
		},
		{
			inpath: ex.f1,
			inRev:  1,
			expErr: false,
			expB:   ex.f1b1,
		},
		{
			inpath: ex.f2,
			inRev:  0,
			expErr: false,
			expB:   ex.f2b1,
		},
		{
			inpath: ex.f2,
			inRev:  5,
			expErr: true,
		},
		{
			inpath: ex.f2 + "randomtext",
			inRev:  0,
			expErr: true,
		},
	}

	for _, c := range cases {
		t.Run("get", func(t *testing.T) {
			f := vcFixture()
			outB, err := f.v.GetFile(c.inpath, c.inRev)
			if c.expErr && err == nil {
				t.Fatalf("expected error")
			}
			if !c.expErr && err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(outB, c.expB) {
				t.Fatalf("wrong get exp %v got %v", c.expB, outB)
			}
		})
	}

}

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
