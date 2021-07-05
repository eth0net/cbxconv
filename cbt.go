package cbxconv

import (
	"archive/tar"
	"fmt"
	"io"
)

type CBTReader struct {
	r *tar.Reader
}

func NewCBTReader(r io.Reader) *CBTReader {
	return &CBTReader{r: tar.NewReader(r)}
}

func (c *CBTReader) Read(p []byte) (int, error) {
	if c.r == nil {
		return 0, fmt.Errorf("nil reader")
	}
	return c.r.Read(p)
}

func (c *CBTReader) Next() (*FileInfo, error) {
	if c.r == nil {
		return nil, fmt.Errorf("no reader")
	}
	hdr, err := c.r.Next()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get next: %s", err)
	}
	fi := &FileInfo{
		Name:    hdr.Name,
		ModTime: hdr.ModTime,
	}
	return fi, nil
}
