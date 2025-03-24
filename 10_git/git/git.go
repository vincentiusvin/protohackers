package git

import (
	"fmt"
	"strings"
)

type VersionControl struct {
	root file
}

func NewVersionControl() *VersionControl {
	return &VersionControl{
		root: newFile(""),
	}
}

// put file
// automatically handle revision
func (v *VersionControl) PutFile(abs_path string, content []byte) (file, error) {
	return nil, nil
}

// get content of file
// revision of 0 means latest revision
func (v *VersionControl) GetFile(abs_path string, revision int) {

}

// list files in a directory
func (v *VersionControl) ListFile(dir string) {

}

func (v *VersionControl) getFile(abs_path string, force bool) (file, error) {
	spls, err := splitPaths(abs_path)
	if err != nil {
		return nil, fmt.Errorf("can't put: %w", err)
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

var errFileName = fmt.Errorf("illegal file name")

func splitPaths(str string) ([]string, error) {
	aft, found := strings.CutPrefix(str, "/")
	if !found || aft == "" {
		return nil, errFileName
	}
	str = aft

	spl := strings.Split(str, "/")
	for _, c := range spl {
		if c == "" {
			return nil, errFileName
		}
	}
	return spl, nil
}

func cutPaths(dir []string) ([]string, string) {
	l := len(dir)
	if l == 0 {
		panic("cutpaths passed a nil string")
	}
	if l == 1 {
		return nil, dir[0]
	}
	return dir[:l-1], dir[l-1]
}
