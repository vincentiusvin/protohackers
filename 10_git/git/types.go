package git

type NodeType int

const (
	nodeDir NodeType = iota + 1
	nodeFile
)

type node interface {
	getName() string
	getType() NodeType
}

type directory struct {
	name    string
	entries []node
}

func newDirectory(name string) *directory {
	return &directory{
		name:    name,
		entries: make([]node, 0),
	}
}

func (d *directory) getName() string {
	return d.name
}

func (d *directory) getType() NodeType {
	return nodeDir
}

func (d *directory) addNode(n node) {
	d.entries = append(d.entries, n)
}

type file struct {
	name string
	// slice of []byte
	// content[0] means the file content of first revision
	content [][]byte
}

func newFile(name string, content []byte) *file {
	f := &file{
		name:    name,
		content: make([][]byte, 0),
	}
	f.addRevision(content)
	return f
}

func (f *file) getName() string {
	return f.name
}

func (f *file) getType() NodeType {
	return nodeFile
}

func (f *file) addRevision(content []byte) {
	f.content = append(f.content, content)
}
