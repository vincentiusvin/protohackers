package git

import "fmt"

type file struct {
	name  string
	child map[string]*file
	// slice of []byte
	// revisions[0] means the file revisions of first revision
	revisions [][]byte
}

func newNode(name string) *file {
	return &file{
		name:      name,
		child:     make(map[string]*file),
		revisions: make([][]byte, 0),
	}
}

func (d *file) getChild(name string) *file {
	return d.child[name]
}

var errNodeExist = fmt.Errorf("file already exists")

func (d *file) addNode(n file) (*file, error) {
	prev := d.child[n.name]
	if prev != nil {
		return nil, errNodeExist
	}
	return prev, nil
}

func (f *file) addRevision(content []byte) int {
	f.revisions = append(f.revisions, content)
	return len(f.revisions)
}
