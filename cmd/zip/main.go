package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	source := filepath.Join("testdata", "comic.cbz")
	r, err := OpenCBZReader(source)
	if err != nil {
		log.Fatalln(err)
	}

	path := filepath.Clean("testdata")
	srcBase := filepath.Base(source)
	srcExt := filepath.Ext(source)
	srcName := strings.TrimSuffix(srcBase, srcExt)
	trgName := srcName + ".out"
	trgExt := ".cbz"
	trgBase := trgName + trgExt
	target := filepath.Join(path, trgBase)

	w, err := OpenCBZWriter(target)
	if err != nil {
		log.Fatalf("failed to create file: %s\n", err)
	}

	for {
		info, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		img, err := jpeg.Decode(r)
		if err != nil {
			log.Fatalf("failed to decode image: %s\n", err)
		}

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, img, nil)
		if err != nil {
			log.Fatalf("failed to encode image: %s\n", err)
		}

		srcExt := filepath.Ext(info.Name)
		trgExt := ".jpg"
		name := strings.TrimSuffix(info.Name, srcExt) + trgExt

		newInfo := &FileInfo{
			Name:    name,
			ModTime: time.Now(),
		}
		_, err = w.WriteWithFileInfo(newInfo, buf.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatalf("failed to close output file: %s\n", err)
	}

	if err = r.Close(); err != nil {
		log.Fatalln(err)
	}
}

type CBZReader struct {
	r  *zip.ReadCloser
	zr io.ReadCloser
	i  int
}

func OpenCBZReader(path string) (*CBZReader, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %s", err)
	}
	return &CBZReader{r: r}, nil
}

func (c *CBZReader) Close() error {
	if c.zr == nil && c.r == nil {
		return fmt.Errorf("nothing to close")
	}
	if c.zr != nil {
		if err := c.zr.Close(); err != nil {
			return fmt.Errorf("failed to close file: %s", err)
		}
	}
	if c.r != nil {
		if err := c.r.Close(); err != nil {
			return fmt.Errorf("failed to close zip: %s", err)
		}
	}
	return nil
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
	f  *os.File
	w  *zip.Writer
	zw io.Writer
}

func OpenCBZWriter(path string) (*CBZWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %s", err)
	}
	return &CBZWriter{f: f, w: zip.NewWriter(f)}, nil
}

func (c *CBZWriter) Close() error {
	if c.f == nil && c.w == nil {
		return fmt.Errorf("nothing to close")
	}
	if c.w != nil {
		if err := c.w.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %s", err)
		}
	}
	if c.f != nil {
		if err := c.f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %s", err)
		}
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

type FileInfo struct {
	Name    string
	ModTime time.Time
}
