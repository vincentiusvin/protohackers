package git

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

type VersionControl struct {
	root file
	mu   sync.Mutex
}

func NewVersionControl() *VersionControl {
	return &VersionControl{
		root: newFile(""),
	}
}

// put file
// automatically handle revision
func (v *VersionControl) PutFile(abs_path string, newData []byte) (int, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	f, err := v.getFile(abs_path, true)
	if err != nil {
		return 0, err
	}
	if f == v.root {
		return 0, errFileName
	}

	if !isText(newData) {
		return 0, errFileContent
	}

	// check if same as prev
	lastRev := f.getRevisionNumber()
	if lastRev != 0 {
		prevData, err := f.getRevision(lastRev)
		if err != nil {
			return 0, errFileNotFound
		}

		if bytes.Equal(prevData, newData) {
			return lastRev, nil
		}
	}

	return f.addRevision(newData), nil
}

// get content of file
// revision of 0 means latest revision
func (v *VersionControl) GetFile(abs_path string, revision int) ([]byte, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	f, err := v.getFile(abs_path, false)
	if err != nil {
		return nil, err
	}
	if f == v.root {
		return nil, errFileName
	}
	rev, err := f.getRevision(revision)
	if err != nil {
		return nil, err
	}
	return rev, nil
}

type FileListItem struct {
	Name string
	Info string
}

// list files in a directory
func (v *VersionControl) ListFile(dir string) ([]FileListItem, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	f, err := v.getFile(dir, false)
	if err != nil {
		return nil, err
	}

	child := f.getChildren()
	ret := make([]FileListItem, len(child))
	for i, f := range child {
		revnum := f.getRevisionNumber()
		ret[i].Name = f.getName()
		if revnum == 0 {
			ret[i].Info = "DIR"
			ret[i].Name += "/"
		} else {
			ret[i].Info = fmt.Sprintf("r%v", revnum)
		}
	}

	slices.SortFunc(ret, func(a, b FileListItem) int {
		if a.Name > b.Name {
			return 1
		} else if b.Name > a.Name {
			return -1
		} else {
			return 0
		}
	})

	return ret, nil
}

func (v *VersionControl) getFile(abs_path string, force bool) (file, error) {
	spls, err := splitPaths(abs_path)
	if err != nil {
		return nil, fmt.Errorf("can't access: %w", err)
	}
	curr := v.root
	for len(spls) != 0 {
		head := spls[0]
		spls = spls[1:]

		f, err := curr.getChild(head)
		if err != nil {
			if err != errFileNotFound || !force {
				return nil, fmt.Errorf("can't get file: %w", err)
			}
			// force && errFileNotFound
			f, err = curr.addChild(head)
			if err != nil {
				return nil, fmt.Errorf("can't create file: %w", err)
			}
		}

		curr = f
	}

	return curr, nil
}

func splitPaths(str string) ([]string, error) {
	aft, found := strings.CutPrefix(str, "/")
	if !found {
		return nil, errFileName
	}
	str = aft
	if str == "" {
		return nil, nil
	}

	spl := strings.Split(str, "/")
	for _, c := range spl {
		if c == "" {
			return nil, errFileName
		}
		if !isAlphanum(c) {
			return nil, errFileName
		}
	}
	return spl, nil
}

func isAlphanum(s string) bool {
	for _, c := range s {
		letter := unicode.IsLetter(c)
		digits := unicode.IsNumber(c)
		dot := c == '.'
		dash := c == '-'
		underscore := c == '_'

		if letter || digits || dot || dash || underscore {
			continue
		}

		return false
	}

	return true
}

// text as defined by the reference implementation
func isText(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}
	s := string(b)
	for _, c := range s {
		ascii := 32 <= c && c <= 127
		letter := unicode.IsLetter(c)
		curr := unicode.Is(unicode.Sc, c)
		space := unicode.IsSpace(c)

		if ascii || letter || curr || space {
			continue
		}
		return false
	}

	return true
}
