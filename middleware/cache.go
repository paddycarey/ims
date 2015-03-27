package middleware

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Cache is a middleware handler that serves cached files from the configured
// filesystem. The URL path and query string are concatenated and hashed
// to form the name of the file that will be looked up from the configured
// filesystem.
type Cache struct {
	// Dir is the directory to serve static files from
	Dir http.FileSystem
}

// NewCache returns a new instance of Cache
func NewCache(directory http.FileSystem) *Cache {
	return &Cache{
		Dir: directory,
	}
}

func (c *Cache) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	cFile := hashURL(r.URL)

	f, err := c.Dir.Open(cFile)
	if err != nil {
		next(rw, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		next(rw, r)
		return
	}

	http.ServeContent(rw, r, r.URL.Path, fi.ModTime(), f)
}

// hashURL concatenates the path and query string of the URL provided and
// hashes the result using md5.
func hashURL(u *url.URL) string {
	hash := md5.New()
	io.WriteString(hash, u.Path)
	io.WriteString(hash, u.RawQuery)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// cacheImage writes content from the provided io.Reader to disk at the
// specified location. Once written files will be available for reading by the
// cache middleware.
func cacheImage(dir string, f io.Reader, name string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	tmpDir := filepath.Join(dir, ".tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	tmpFile, err := ioutil.TempFile(tmpDir, "")
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	tmpFile.Write(b)
	if err := tmpFile.Sync(); err != nil {
		return err
	}

	tmpFile.Close()
	if err := os.Rename(tmpFile.Name(), filepath.Join(dir, name)); err != nil {
		os.Remove(tmpFile.Name())
		return err
	}

	return nil
}
