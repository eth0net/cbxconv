package cbxconv

import (
	"io"
	"time"
)

type Reader interface {
	io.Reader
	Next() (fi *FileInfo, err error)
}

type Writer interface {
	io.Closer
	io.Writer
	WriteFileInfo(i *FileInfo) error
	WriteWithFileInfo(i *FileInfo, p []byte) (n int, err error)
}

type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}
