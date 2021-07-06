package main

import (
	"bytes"
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
	source := filepath.Join("testdata", "single.jpg.cbz")

	fi, err := os.Open(source)
	if err != nil {
		log.Fatalf("failed to open input archive: %s\n", err)
	}
	ii, err := fi.Stat()
	if err != nil {
		log.Fatalf("failed to stat input archive: %s\n", err)
	}

	path := filepath.Clean("testdata")
	srcBase := filepath.Base(source)
	srcExt := filepath.Ext(source)
	srcName := strings.TrimSuffix(srcBase, srcExt)
	trgName := srcName + ".out"
	trgExt := ".cbz"
	trgBase := trgName + trgExt
	target := filepath.Join(path, trgBase)

	fo, err := os.Create(target)
	if err != nil {
		log.Fatalf("failed to create file: %s\n", err)
	}

	r, err := cbxconv.NewCBZReader(fi, ii.Size())
	if err != nil {
		log.Fatalf("failed to create archive reader: %s\n", err)
	}
	w, err := cbxconv.NewCBZWriter(fo)
	if err != nil {
		log.Fatalf("failed to create archive writer: %s\n", err)
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

		newInfo := &cbxconv.FileInfo{
			Name:    name,
			ModTime: time.Now(),
		}
		_, err = w.WriteWithFileInfo(newInfo, buf.Bytes())
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err = w.Close(); err != nil {
		log.Fatalf("failed to close output archive: %s\n", err)
	}

	if err = fo.Close(); err != nil {
		log.Fatalf("failed to close output file: %s\n", err)
	}

	if err = fi.Close(); err != nil {
		log.Fatalln(err)
	}
}
