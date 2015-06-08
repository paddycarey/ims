package storage

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type LocalFile struct {
	*os.File
}

func (lf *LocalFile) MimeType() string {
	fi, err := lf.Stat()
	if err != nil {
		return ""
	}
	return getMimeTypeFromFilename(fi.Name())
}

func (lf *LocalFile) ModTime() time.Time {
	fi, err := lf.Stat()
	if err != nil {
		return time.Time{}
	}
	return fi.ModTime()
}

type LocalFileSystem struct {
	dir string
}

func (l *LocalFileSystem) Open(name string) (File, error) {
	if filepath.Separator != '/' && strings.IndexRune(name, filepath.Separator) >= 0 ||
		strings.Contains(name, "\x00") {
		return nil, errors.New("invalid character in file path")
	}
	dir := string(l.dir)
	if dir == "" {
		dir = "."
	}
	f, err := os.Open(filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))))
	if err != nil {
		return nil, err
	}
	lf := &LocalFile{f}
	return lf, nil
}

func NewLocalFileSystem(dir string) (*LocalFileSystem, error) {

	src, err := os.Stat(dir)
	if err != nil {
		return &LocalFileSystem{}, err
	}

	if !src.IsDir() {
		fmt.Errorf("%s is not a directory", dir)
		return &LocalFileSystem{}, err
	}

	return &LocalFileSystem{dir}, nil
}
