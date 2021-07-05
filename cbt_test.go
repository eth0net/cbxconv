package cbxconv

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type cbtTest struct {
	name  string
	files []cbtTestFile
	err   error
}

type cbtTestFile struct {
	name   string
	bytes  []byte
	ignore bool
}

var cbtTests = []cbtTest{
	{
		name: "multi.info.jpg.cbt",
		files: []cbtTestFile{
			{name: "1.jpg"},
			{name: "2.jpg"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "multi.info.png.cbt",
		files: []cbtTestFile{
			{name: "1.png"},
			{name: "2.png"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "multi.info.webp.cbt",
		files: []cbtTestFile{
			{name: "1.webp"},
			{name: "2.webp"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "multi.jpg.cbt",
		files: []cbtTestFile{
			{name: "1.jpg"},
			{name: "2.jpg"},
		},
	},
	{
		name: "multi.png.cbt",
		files: []cbtTestFile{
			{name: "1.png"},
			{name: "2.png"},
		},
	},
	{
		name: "multi.webp.cbt",
		files: []cbtTestFile{
			{name: "1.webp"},
			{name: "2.webp"},
		},
	},
	{
		name: "single.info.jpg.cbt",
		files: []cbtTestFile{
			{name: "1.jpg"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "single.info.png.cbt",
		files: []cbtTestFile{
			{name: "1.png"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "single.info.webp.cbt",
		files: []cbtTestFile{
			{name: "1.webp"},
			{name: "ComicInfo.xml", ignore: true},
		},
	},
	{
		name: "single.jpg.cbt",
		files: []cbtTestFile{
			{name: "1.jpg"},
		},
	},
	{
		name: "single.png.cbt",
		files: []cbtTestFile{
			{name: "1.png"},
		},
	},
	{
		name: "single.webp.cbt",
		files: []cbtTestFile{
			{name: "1.webp"},
		},
	},
}

func TestNewCBTReader(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	want := &CBTReader{r: tar.NewReader(buf)}
	got := NewCBTReader(buf)
	if !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %#v, got %#v", want, got)
	}
}

func TestCBTReader_Read(t *testing.T) {
	t.Parallel()

	t.Run("nil reader", func(t *testing.T) {
		r := CBTReader{}
		_, err := r.Read([]byte{})
		want := fmt.Errorf("nil reader")
		if !reflect.DeepEqual(want, err) {
			t.Errorf("wanted err == %q, got %q", want, err)
		}
	})

	for _, c := range cbtTests {
		t.Run(c.name, func(t *testing.T) {
			f, err := os.Open(filepath.Join("testdata", c.name))
			if err != nil {
				t.Errorf("failed to open %q", c.name)
				return
			}
			r := NewCBTReader(f)

			for _, file := range c.files {
				if _, err = r.Next(); err != nil {
					if err != io.EOF {
						t.Errorf("unexpected error: %q", err)
					}
					return
				}

				var want, got []byte
				got, err = io.ReadAll(r)
				if err != c.err && err != io.EOF {
					t.Errorf("wanted err == %q, got %q", c.err, err)
					break
				}

				if file.ignore {
					continue
				}

				var f1 *os.File
				f1, err = os.Open(filepath.Join("testdata", file.name))
				if err != nil {
					t.Errorf("failed to open %q", file.name)
					return
				}
				want, err = io.ReadAll(f1)
				if err != nil {
					t.Errorf("failed to read %q", file.name)
					return
				}
				if err = f1.Close(); err != nil {
					t.Errorf("failed to close %q", file.name)
					return
				}

				if !reflect.DeepEqual(want, got) {
					t.Errorf("%q does not match", file.name)
				}
			}

			if err = f.Close(); err != nil {
				t.Errorf("failed to close %q", c.name)
			}
		})
	}
}

func TestCBTReader_Next(t *testing.T) {

}
