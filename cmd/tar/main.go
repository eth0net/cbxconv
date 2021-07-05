package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/eth0net/cbxconv"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	source := filepath.Join("testdata", "128.cbt")

	fi, err := os.Open(source)
	if err != nil {
		log.Fatalf("failed to open input file: %s", err)
	}

	path := filepath.Clean("testdata")
	srcBase := filepath.Base(source)
	srcExt := filepath.Ext(source)
	srcName := strings.TrimSuffix(srcBase, srcExt)
	trgName := srcName + ".out"
	trgExt := ".cbt"
	trgBase := trgName + trgExt
	target := filepath.Join(path, trgBase)

	fo, err := os.Create(target)
	if err != nil {
		log.Fatalf("failed to create output file: %s", err)
	}

	r := cbxconv.NewCBTReader(fi)
	w := NewCBTWriter(fo)

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

		newInfo := &cbxconv.FileInfo{
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

	if err := fo.Close(); err != nil {
		log.Fatalln(err)
	}
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

func (c *CBTWriter) WriteFileInfo(i *cbxconv.FileInfo) error {
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

func (c *CBTWriter) WriteWithFileInfo(i *cbxconv.FileInfo, p []byte) (int, error) {
	if err := c.WriteFileInfo(i); err != nil {
		return 0, err
	}
	return c.Write(p)
}
