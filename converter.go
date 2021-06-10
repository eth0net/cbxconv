package cbxconv

import (
	"fmt"
	"os"
	"sync"
)

// A Converter takes comic book archives, converts them to the
// configured formats, and then writes them to the output path.
type Converter struct {
	mu                   sync.Mutex
	archive, image, path string
}

// New creates a new Converter with the provided parameters.
func New(archive, image, path string) (*Converter, error) {
	c := &Converter{}
	err := c.SetArchive(archive)
	if err != nil {
		return nil, err
	}
	err = c.SetImage(image)
	if err != nil {
		return nil, err
	}
	err = c.SetPath(path)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// SetArchive sets the archive format used by the Converter.
// An error will be returned if archive is not a supported archive format.
func (c *Converter) SetArchive(archive string) error {
	switch archive {
	case ArchiveRar, ArchiveTar, ArchiveZip:
		c.mu.Lock()
		c.archive = archive
		c.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("unknown archive format: %q", archive)
	}
}

// SetImage sets the image format used by the Converter.
// An error will be returned if image is not a supported image format.
func (c *Converter) SetImage(image string) error {
	switch image {
	case ImageJPG, ImagePNG, ImageWEBP:
		c.mu.Lock()
		c.image = image
		c.mu.Unlock()
		return nil
	default:
		return fmt.Errorf("unknown image format: %q", image)
	}
}

// SetPath sets the path format used by the Converter.
// An error will be returned if path is not a valid output path.
func (c *Converter) SetPath(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %q", path)
	}
	if !f.IsDir() {
		return fmt.Errorf("path is not a directory: %q", path)
	}
	c.mu.Lock()
	c.path = path
	c.mu.Unlock()
	return nil
}
