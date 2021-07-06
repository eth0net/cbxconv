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
		return nil, fmt.Errorf("nil reader")
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

type CBTWriter struct {
	w *tar.Writer
}

func NewCBTWriter(w io.Writer) *CBTWriter {
	return &CBTWriter{w: tar.NewWriter(w)}
}

func (c *CBTWriter) Close() error {
	if c.w == nil {
		return fmt.Errorf("nothing to close")
	}
	if err := c.w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %s", err)
	}
	return nil
}

func (c *CBTWriter) Write(p []byte) (int, error) {
	if c.w == nil {
		return 0, fmt.Errorf("no writer")
	}
	n, err := c.w.Write(p)
	if err != nil {
		err = fmt.Errorf("failed to write bytes: %s", err)
	}
	return n, err
}

func (c *CBTWriter) WriteFileInfo(i *FileInfo) error {
	if c.w == nil {
		return fmt.Errorf("no writer")
	}
	hdr := &tar.Header{
		Name:    i.Name,
		Size:    i.Size,
		Mode:    0600,
		ModTime: i.ModTime,
	}
	err := c.w.WriteHeader(hdr)
	if err != nil {
		return fmt.Errorf("failed to write header: %s", err)
	}
	return nil
}

func (c *CBTWriter) WriteWithFileInfo(i *FileInfo, p []byte) (int, error) {
	if err := c.WriteFileInfo(i); err != nil {
		return 0, err
	}
	return c.Write(p)
}
