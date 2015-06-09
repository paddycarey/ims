package cache

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	gocache "github.com/pmylund/go-cache"

	"github.com/paddycarey/ims/pkg/storage"
)

type InMemoryCache struct {
	c *gocache.Cache
}

func NewInMemoryCache() (*InMemoryCache, error) {
	return &InMemoryCache{
		c: gocache.New(60*time.Minute, 30*time.Second),
	}, nil
}

// GenerateKey concatenates the path and query string of the URL provided and
// hashes the result using md5.
func (c *InMemoryCache) GenerateKey(r *http.Request) string {
	hash := md5.New()
	io.WriteString(hash, r.URL.Path)
	io.WriteString(hash, r.URL.RawQuery)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (c *InMemoryCache) Get(ck string) (storage.File, bool, error) {
	file, found := c.c.Get(ck)
	if found {
		return file.(storage.File), true, nil
	}
	return nil, false, errors.New("File not found in cache")
}

func (c *InMemoryCache) Set(ck string, f storage.File) error {
	c.c.Set(ck, f, gocache.DefaultExpiration)
	return nil
}
