package server

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/paddycarey/ims/filters"
	"github.com/paddycarey/ims/storage"
)

// FilterMiddleware is a middleware handler that performs image manipulation on
// the fly, caching the results in-memory.
type FilterMiddleware struct {
	// storageBackend can be anything that implements the `http.FileSystem`
	// interface. This allows easy use of either a local directory or a remote
	// store like S3 or GCS.
	storageBackend http.FileSystem
	cacheDir       string
}

// NewFilterMiddleware returns a new instance of FilterMiddleware
func NewFilterMiddleware(storageBackend http.FileSystem, cacheDir string) *FilterMiddleware {
	return &FilterMiddleware{
		storageBackend: storageBackend,
		cacheDir:       cacheDir,
	}
}

func (f *FilterMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	var serveFile io.ReadSeeker

	// check that the given path ends with a known image extension.
	format, err := validateImageExtension(r.URL.Path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Invalid image extension")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("400 Bad Request"))
		return
	}

	// load file from storage backend
	file, err := f.storageBackend.Open(r.URL.Path)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Unable to open file")
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("404 Not Found"))
		return
	}
	defer file.Close()
	serveFile = file

	stat, err := file.Stat()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Unable to stat file")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("500 Internal Server Error"))
		return
	}
	modTime := stat.ModTime()

	if r.URL.RawQuery != "" {

		img, _, err := storage.DecodeImage(file)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Unable to parse file as image")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("500 Internal Server Error"))
			return
		}

		g, err := filters.ParseFilters(r.URL.RawQuery)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Unable to parse filters")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("400 Bad Request"))
			return
		}
		err = img.ApplyFilters(g)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Unable to apply filters")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("500 Internal Server Error"))
			return
		}
		modTime = time.Now()

		serveFile, err = img.Encode()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("Unable to encode image")
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("500 Internal Server Error"))
			return
		}

	}

	// store the image in the cache directory
	if err := cacheImage(serveFile, f.cacheDir, hashURL(r.URL)); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("Unable to cache image")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("500 Internal Server Error"))
		return
	}

	rw.Header().Set("Content-Type", fmt.Sprintf("image/%s", format))
	http.ServeContent(rw, r, r.URL.Path, modTime, serveFile)
}

// validateImageExtension checks if the given string ends in a known
// valid/supported image extension.
func validateImageExtension(path string) (string, error) {

	extension := filepath.Ext(path)
	extension = strings.ToLower(extension)
	switch extension {
	case ".gif":
		return "gif", nil
	case ".jpeg", ".jpg":
		return "jpeg", nil
	case ".png":
		return "png", nil
	}

	return "", fmt.Errorf("Matching extension not found: %s", path)
}
