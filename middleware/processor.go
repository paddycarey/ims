package middleware

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/paddycarey/ims/processors"
)

// Processor is a middleware handler that performs image manipulation on the
// fly, writing the results to durable storage for caching purposes.
type Processor struct {
	// storageBackend can be anything that implements the `http.FileSystem`
	// interface. This allows easy use of either a local directory or a remote
	// store like S3 or GCS.
	storageBackend http.FileSystem
	cacheBackend   string
}

// NewProcessor returns a new instance of Processor
func NewProcessor(storageBackend http.FileSystem, cacheBackend string) *Processor {
	return &Processor{
		storageBackend: storageBackend,
		cacheBackend:   cacheBackend,
	}
}

func (p *Processor) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// check that the given path ends with a known image extension.
	if !validateImageExtension(r.URL.Path) {
		next(rw, r)
		return
	}

	// load file from storage backend, falling through if unable
	f, err := p.storageBackend.Open(r.URL.Path)
	if err != nil {
		next(rw, r)
		return
	}
	defer f.Close()

	// serve file without any processing if no arguments were passed in the
	// query string.
	if r.URL.RawQuery == "" {
		stat, err := f.Stat()
		if err != nil {
			next(rw, r)
			return
		}
		http.ServeContent(rw, r, r.URL.Path, stat.ModTime(), f)
		return
	}

	// process the image according to the arguments passed in the url
	pf, err := processors.ProcessImage(r, f)
	if err != nil {
		next(rw, r)
		return
	}

	// store the image in the cache directory
	if err := cacheImage(p.cacheBackend, pf, hashURL(r.URL)); err != nil {
		log.Printf("%s", err)
	}

	http.ServeContent(rw, r, r.URL.Path, time.Time{}, pf)
}

// validateImageExtension checks if the given string ends in a known
// valid/supported image extension.
func validateImageExtension(path string) bool {

	extension := filepath.Ext(path)
	extension = strings.ToLower(extension)
	switch extension {
	case ".gif", ".jpeg", ".jpg", ".png":
		return true
	}

	return false
}
