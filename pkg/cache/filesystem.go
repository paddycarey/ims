package cache

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/paddycarey/ims/pkg/storage"
)

type LocalFileSystemCache struct {
	dir string
	fs  storage.FileSystem
}

func NewLocalFileSystemCache(dir string) (*LocalFileSystemCache, error) {
	// ensure directory exists and is writable
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// init a new FileSystem we can use to load files from disk
	fs, err := storage.NewLocalFileSystem(dir)
	if err != nil {
		return nil, err
	}

	return &LocalFileSystemCache{dir: dir, fs: fs}, nil
}

// GenerateKey concatenates the path and query string of the URL provided and
// hashes the result using md5.
func (c *LocalFileSystemCache) GenerateKey(r *http.Request) string {
	hash := md5.New()
	io.WriteString(hash, r.URL.Path)
	io.WriteString(hash, r.URL.RawQuery)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (c *LocalFileSystemCache) Get(ck string) (storage.File, bool, error) {
	f, err := c.fs.Open(ck)
	if err != nil {
		return nil, false, err
	}
	return f, true, nil
}

func (c *LocalFileSystemCache) Set(ck string, f storage.File) error {
	// ensure temporary directory exists and is writable
	tmpDir := filepath.Join(c.dir, ".tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}

	// create a temporary file we can write to before renaming the newly
	// written file into its final place
	tmpFile, err := ioutil.TempFile(tmpDir, "")
	if err != nil {
		return err
	}

	// write data to temporary file, explicitly syncing and closing it
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tmpFile.Write(b)
	if err := tmpFile.Sync(); err != nil {
		return err
	}
	tmpFile.Close()

	if err := os.Rename(tmpFile.Name(), filepath.Join(c.dir, ck)); err != nil {
		os.Remove(tmpFile.Name())
		return err
	}
	return nil
}
