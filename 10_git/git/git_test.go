package git_test

import (
	"bytes"
	"protohackers/10_git/git"
	"reflect"
	"testing"
)

type fixtureRet struct {
	v    *git.VersionControl
	f1   string
	f2   string
	f3   string
	f4   string
	f1b1 []byte
	f1b2 []byte
	f2b1 []byte
	f3b1 []byte
	f4b1 []byte
}

func vcFixture() fixtureRet {
	v := git.NewVersionControl()
	f1b1 := []byte{0x01, 0x02}
	f1b2 := []byte{0x01, 0x04}
	f2b1 := []byte{0x01, 0x03}
	f3b1 := []byte{0x01, 0x07}
	f4b1 := []byte{0x01, 0x08}
	f1 := "/dir1/dirfile/file"
	f2 := "/dir1/dirfile"
	f3 := "/dir1/dir/file2"
	f4 := "/dir1/dir/file1"
	v.PutFile(f1, f1b1)
	v.PutFile(f1, f1b2)
	v.PutFile(f2, f2b1)
	v.PutFile(f3, f3b1)
	v.PutFile(f4, f4b1)
	return fixtureRet{
		v:    v,
		f1:   f1,
		f2:   f2,
		f3:   f3,
		f4:   f4,
		f1b1: f1b1,
		f1b2: f1b2,
		f2b1: f2b1,
		f3b1: f3b1,
		f4b1: f4b1,
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

func TestList(t *testing.T) {
	f := vcFixture()

	type listCases struct {
		in      string
		entries []git.FileListItem
	}

	cases := []listCases{
		{
			in: "/dir1",
			entries: []git.FileListItem{
				{
					Name: "dir/",
					Info: "DIR",
				},
				{
					Name: "dirfile",
					Info: "r1",
				},
			},
		},
		{
			in: "/",
			entries: []git.FileListItem{
				{
					Name: "dir1/",
					Info: "DIR",
				},
			},
		},
		{
			in: "/dir1/dir",
			entries: []git.FileListItem{
				{
					Name: "file1",
					Info: "r1",
				},
				{
					Name: "file2",
					Info: "r1",
				},
			},
		},
	}

	for _, c := range cases {
		t.Run("case", func(t *testing.T) {
			ret, _ := f.v.ListFile(c.in)
			if !reflect.DeepEqual(ret, c.entries) {
				t.Fatalf("wrong entires exp %v got %v", c.entries, ret)
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
