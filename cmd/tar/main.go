package main

import (
	"archive/tar"
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
	source := filepath.Join("testdata", "comic.cbt")
	r, err := OpenCBTReader(source)
	if err != nil {
		log.Fatalln(err)
	}

	path := filepath.Clean("testdata")
	srcBase := filepath.Base(source)
	srcExt := filepath.Ext(source)
	srcName := strings.TrimSuffix(srcBase, srcExt)
	trgName := srcName + ".out"
	trgExt := ".cbt"
	trgBase := trgName + trgExt
	target := filepath.Join(path, trgBase)

	w, err := OpenCBTWriter(target)
	if err != nil {
		log.Fatalln(err)
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
			Size:    int64(buf.Len()),
		}
		_, err = w.WriteWithFileInfo(newInfo, buf.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := w.Close(); err != nil {
		log.Fatalln(err)
	}

	if err := r.Close(); err != nil {
		log.Fatalln(err)
	}
}

type CBTReader struct {
	f *os.File
	r *tar.Reader
}

func OpenCBTReader(path string) (*CBTReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %s", err)
	}
	return &CBTReader{f: f, r: tar.NewReader(f)}, nil
}

func (c *CBTReader) Close() error {
	if c.f == nil {
		return fmt.Errorf("nothing to close")
	}
	if err := c.f.Close(); err != nil {
		return fmt.Errorf("failed to close: %s", err)
	}
	return nil
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

func (c *CBTReader) Read(p []byte) (int, error) {
	if c.r == nil {
		return 0, fmt.Errorf("no reader")
	}
	n, err := c.r.Read(p)
	if err != nil {
		return n, fmt.Errorf("failed to read: %s", err)
	}
	return n, err
}

type CBTWriter struct {
	f *os.File
	w *tar.Writer
}

func OpenCBTWriter(path string) (*CBTWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %s", err)
	}
	return &CBTWriter{f: f, w: tar.NewWriter(f)}, nil
}

func (c *CBTWriter) Close() error {
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

type FileInfo struct {
	Name    string
	ModTime time.Time
	Size    int64
}
