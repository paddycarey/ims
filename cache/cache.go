package cache

import (
	"net/http"

	"github.com/paddycarey/ims/storage"
)

type CacheBackend interface {
	GenerateKey(*http.Request) string
	Get(string) (storage.File, bool, error)
	Set(string, storage.File) error
}

func LoadBackend(uri string) (CacheBackend, error) {
	// fall through to default filesystem loader
	return NewLocalFileSystemCache(uri)
}
