package git

import (
	"fmt"
	"strings"
	"sync"
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
func (v *VersionControl) PutFile(abs_path string, content []byte) (int, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	f, err := v.getFile(abs_path, true)
	if err != nil {
		return 0, err
	}
	revnum := f.addRevision(content)
	return revnum, nil
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

	return ret, nil
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
