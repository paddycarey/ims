package server

import (
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/paddycarey/ims/pkg/images"
	"github.com/paddycarey/ims/pkg/storage"
)

type Server struct {
	Cache   *InMemoryCache
	Storage storage.FileSystem
	// disable optimizations
	NoOpts bool
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	// generate a cache key from the incoming request
	ck := s.Cache.GenerateKey(r)

	// attempt to load the processed image from cache
	processedImage, found, _ := s.Cache.Get(ck)
	if found {
		s.serveImage(rw, r, processedImage)
		return
	}

	// load image from storage backend
	storedImage, err := s.Storage.Open(r.URL.Path)
	if err != nil {
		s.serveError(rw, r, err)
		return
	}
	defer storedImage.Close()

	// process the image that's been loaded from storage
	processedImage, err = images.Process(storedImage, r.URL, s.NoOpts)
	if err != nil {
		s.serveError(rw, r, err)
		return
	}

	// cache the newly processed image
	s.Cache.Set(ck, processedImage)

	// serve the newly processed image
	s.serveImage(rw, r, processedImage)
}

func (s *Server) serveError(rw http.ResponseWriter, r *http.Request, err error) {

	var status string

	if strings.Contains(err.Error(), "no such file or directory") {
		rw.WriteHeader(http.StatusNotFound)
		status = "404 Not Found"
	} else if strings.Contains(err.Error(), "404") {
		rw.WriteHeader(http.StatusNotFound)
		status = "404 Not Found"
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		status = "500 Internal Server Error"
	}

	rw.Write([]byte(status))
	logrus.WithFields(logrus.Fields{
		"err":    err,
		"status": status,
		"url":    r.URL.String(),
	}).Error("Error serving request")
}

func (s *Server) serveImage(rw http.ResponseWriter, r *http.Request, f storage.File) {
	rw.Header().Set("Content-Type", f.MimeType())
	http.ServeContent(rw, r, r.URL.Path, f.ModTime(), f)
}
