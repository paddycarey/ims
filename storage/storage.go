package storage

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type File interface {
	io.Closer
	io.ReadSeeker
	MimeType() string
	ModTime() time.Time
}

type FileSystem interface {
	Open(string) (File, error)
}

func LoadBackend(uri string) (FileSystem, error) {

	if strings.HasPrefix(uri, "s3://") {
		err := errors.New("s3 loader not implemented yet")
		return nil, err
	}

	if strings.HasPrefix(uri, "gcs://") {
		return NewGCSFileSystem(uri)
	}

	// fall through to default filesystem loader
	return NewLocalFileSystem(uri)
}

func getMimeTypeFromFilename(name string) string {
	ext := filepath.Ext(name)
	switch ext {
	case ".gif":
		return "image/gif"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	}

	return ""
}
