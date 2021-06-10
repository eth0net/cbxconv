package cbxconv

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		archive, image, path string
		want                 *Converter
		wantErr              bool
	}{
		{
			name:    "rar with jpg in current dir",
			archive: ArchiveRar,
			image:   ImageJPG,
			path:    ".",
			want: &Converter{
				archive: ArchiveRar,
				image:   ImageJPG,
				path:    ".",
			},
			wantErr: false,
		},
		{
			name:    "tar with png in current dir",
			archive: ArchiveTar,
			image:   ImagePNG,
			path:    ".",
			want: &Converter{
				archive: ArchiveTar,
				image:   ImagePNG,
				path:    ".",
			},
			wantErr: false,
		},
		{
			name:    "zip with webp in current dir",
			archive: ArchiveZip,
			image:   ImageWEBP,
			path:    ".",
			want: &Converter{
				archive: ArchiveZip,
				image:   ImageWEBP,
				path:    ".",
			},
			wantErr: false,
		},
		{
			name:    "bad archive",
			archive: "bad",
			image:   ImageWEBP,
			path:    ".",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "bad image",
			archive: ArchiveZip,
			image:   "bad",
			path:    ".",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "bad dir",
			archive: ArchiveZip,
			image:   ImageWEBP,
			path:    "bad",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		got, err := New(tc.archive, tc.image, tc.path)
		if (err != nil) != tc.wantErr {
			t.Fatalf("case %q: wantErr = %t, got err = %+v", tc.name, tc.wantErr, err)
		}
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("case %q: want %+v, got %+v", tc.name, tc.want, got)
		}
	}
}

func TestConverter_SetArchive(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input, want string
		wantErr     bool
	}{
		{"rar archive", ArchiveRar, ArchiveRar, false},
		{"tar archive", ArchiveTar, ArchiveTar, false},
		{"zip archive", ArchiveZip, ArchiveZip, false},
		{"bad archive", "bad", "", true},
	}

	for _, tc := range testCases {
		c := Converter{}
		err := c.SetArchive(tc.input)
		if (err != nil) != tc.wantErr {
			t.Fatalf("case %q: wantErr = %t, got err = %+v", tc.name, tc.wantErr, err)
		}
		got := c.archive
		if tc.want != got {
			t.Errorf("case %q: want %q, got %q", tc.name, tc.want, got)
		}
	}
}

func TestConverter_SetImage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input, want string
		wantErr     bool
	}{
		{"jpg image", ImageJPG, ImageJPG, false},
		{"png image", ImagePNG, ImagePNG, false},
		{"webp image", ImageWEBP, ImageWEBP, false},
		{"bad image", "bad", "", true},
	}

	for _, tc := range testCases {
		c := Converter{}
		err := c.SetImage(tc.input)
		if (err != nil) != tc.wantErr {
			t.Fatalf("case %q: wantErr = %t, got err = %+v", tc.name, tc.wantErr, err)
		}
		got := c.image
		if tc.want != got {
			t.Errorf("case %q: want %q, got %q", tc.name, tc.want, got)
		}
	}
}

func TestConverter_SetPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		input, want string
		wantErr     bool
	}{
		{"current dir is valid", ".", ".", false},
		{"non-existent dir", "bad", "", true},
		{"not a dir", "go.mod", "", true},
	}

	for _, tc := range testCases {
		c := Converter{}
		err := c.SetPath(tc.input)
		if (err != nil) != tc.wantErr {
			t.Fatalf("case %q: wantErr = %t, got err = %+v", tc.name, tc.wantErr, err)
		}
		got := c.path
		if tc.want != got {
			t.Errorf("case %q: want %q, got %q", tc.name, tc.want, got)
		}
	}
}
