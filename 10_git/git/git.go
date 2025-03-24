package git

import (
	"fmt"
	"strings"
)

type VersionControl struct {
	entries map[string]node
}

func NewVersionControl() *VersionControl {
	return &VersionControl{
		entries: make(map[string]node),
	}
}

// put file
// automatically handle revision
func (v *VersionControl) PutFile(abs_path string, content []byte) {

}

// get content of file
// revision of 0 means latest revision
func (v *VersionControl) GetFile(abs_path string, revision int) {

}

// list files in a directory
func (v *VersionControl) ListFile(dir string) {

}

var errFileName = fmt.Errorf("illegal file name")

func splitPaths(str string) ([]string, error) {
	aft, found := strings.CutPrefix(str, "/")
	if !found {
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
