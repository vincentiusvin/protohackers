package git

import "fmt"

var (
	errRevNotFound  = fmt.Errorf("revision not found")
	errFileNotFound = fmt.Errorf("file not found")
	errNodeExist    = fmt.Errorf("file already exists")
	errFileName     = fmt.Errorf("illegal file name")
	errFileContent  = fmt.Errorf("illegal file content")
)

// readonly to the outside world
type file interface {
	getName() string
	getRevisionNumber() (revnum int)
	getChildren() (f []file)
	getChild(name string) (f file, err error)
	addChild(name string) (f file, err error)
	getRevision(rev int) (data []byte, err error)
	addRevision(data []byte) (revnum int)
}

type cfile struct {
	name  string
	child map[string]file
	// slice of []byte
	// revisions[0] means the file revisions of first revision
	revisions [][]byte
}

func newFile(name string) file {
	f := &cfile{
		name:      name,
		child:     make(map[string]file),
		revisions: make([][]byte, 0),
	}
	return f
}

func (d *cfile) getName() string {
	return d.name
}

func (d *cfile) getRevisionNumber() int {
	return len(d.revisions)
}

func (d *cfile) getChild(name string) (file, error) {
	f, ok := d.child[name]
	if !ok {
		return nil, errFileNotFound
	}
	return f, nil
}

func (d *cfile) addChild(name string) (file, error) {
	prev, _ := d.getChild(name)
	if prev != nil {
		return nil, errNodeExist
	}
	d.child[name] = newFile(name)
	return d.child[name], nil
}

func (f *cfile) getRevision(rev int) ([]byte, error) {
	l := len(f.revisions)
	if l == 0 {
		return nil, errFileNotFound
	}
	if rev < 0 || rev > l {
		return nil, errRevNotFound
	}
	if rev == 0 {
		return f.revisions[l-1], nil
	}
	return f.revisions[rev-1], nil
}

func (f *cfile) addRevision(data []byte) int {
	f.revisions = append(f.revisions, data)
	return len(f.revisions)
}

func (f *cfile) getChildren() []file {
	ret := make([]file, 0)
	for _, fc := range f.child {
		ret = append(ret, fc)
	}
	return ret
}
