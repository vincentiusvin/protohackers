package git

import (
	"log"
	"testing"
)

func TestDirectory(t *testing.T) {
	dir := newDirectory("dir1")
	n := newFile("file1", []byte{0x01, 0x02})
	dir.addNode(n)
	log.Println(n)
}
