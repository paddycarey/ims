package storage

import (
	"bytes"
	"time"
)

type InMemoryFile struct {
	*bytes.Reader
	Mimetype string
	Modtime  time.Time
}

func (i *InMemoryFile) Close() error {
	return nil
}

func (i *InMemoryFile) MimeType() string {
	return i.Mimetype
}

func (i *InMemoryFile) ModTime() time.Time {
	return i.Modtime
}
