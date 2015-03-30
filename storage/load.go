package storage

import (
	"net/http"
)

func LoadBackend(uri string) (http.Dir, error) {
	return NewFileSystemStorage(uri)
}
