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
	source := filepath.Join("testdata", "single.jpg.cbt")

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
	w := cbxconv.NewCBTWriter(fo)

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
