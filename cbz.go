package cbxconv

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
)

type CBZReader struct {
	r  *zip.Reader
	zr io.ReadCloser
	i  int
}

func NewCBZReader(r io.ReaderAt, size int64) (*CBZReader, error) {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	return &CBZReader{r: zr}, nil
}

func (c *CBZReader) Next() (fi *FileInfo, err error) {
	if c.r == nil {
		return nil, fmt.Errorf("no reader")
	}
	if c.zr != nil {
		if err := c.zr.Close(); err != nil {
			return nil, fmt.Errorf("failed to close current file: %s", err)
		}
		c.zr = nil
	}
	if c.i >= len(c.r.File) {
		return nil, io.EOF
	}
	f := c.r.File[c.i]
	c.zr, err = f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open next file: %s", err)
	}
	c.i++
	fi = &FileInfo{
		Name:    f.Name,
		ModTime: f.Modified,
	}
	return fi, nil
}

func (c *CBZReader) Read(p []byte) (int, error) {
	if c.zr == nil {
		return 0, fmt.Errorf("no reader")
	}
	n, err := c.zr.Read(p)
	if err != nil {
		return n, fmt.Errorf("failed to read: %s", err)
	}
	return n, nil
}

type CBZWriter struct {
	w  *zip.Writer
	zw io.Writer
}

func NewCBZWriter(w io.Writer) (*CBZWriter, error) {
	return &CBZWriter{w: zip.NewWriter(w)}, nil
}

func (c *CBZWriter) Close() error {
	if c.w == nil {
		return fmt.Errorf("nothing to close")
	}
	if err := c.w.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %s", err)
	}
	return nil
}

func (c *CBZWriter) Write(p []byte) (int, error) {
	if c.zw == nil {
		return 0, fmt.Errorf("no writer")
	}
	n, err := c.zw.Write(p)
	if err != nil {
		err = fmt.Errorf("failed to write bytes: %s", err)
	}
	return n, err
}

func (c *CBZWriter) WriteFileInfo(i *FileInfo) error {
	if c.w == nil {
		return fmt.Errorf("no writer")
	}
	hdr := &zip.FileHeader{
		Name:     i.Name,
		Modified: i.ModTime,
	}
	w, err := c.w.CreateHeader(hdr)
	if err != nil {
		log.Fatalf("failed to write header: %s\n", err)
	}
	c.zw = w
	return nil
}

func (c *CBZWriter) WriteWithFileInfo(i *FileInfo, p []byte) (int, error) {
	err := c.WriteFileInfo(i)
	if err != nil {
		return 0, err
	}
	n, err := c.Write(p)
	return n, err
}
