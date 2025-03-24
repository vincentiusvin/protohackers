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

func (v *VersionControl) getFile(abs_path string) (file, error) {
	_, err := splitPaths(abs_path)
	if err != nil {
		return nil, fmt.Errorf("can't put: %w", err)
	}

	return nil, nil
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

func cutPaths([]string) (dir []string, filename string) {
	l := len(dir)
	if l == 0 {
		panic("cutpaths passed a nil string")
	}
	if l == 1 {
		return nil, dir[0]
	}
	return dir[:l-1], dir[l-1]
}
