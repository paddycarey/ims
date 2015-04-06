package storage

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func LoadBackend(uri string) (http.Dir, error) {
	missingLoaders := []string{"s3", "gcs"}
	for _, loader := range missingLoaders {
		if strings.HasPrefix(uri, loader) {
			err := fmt.Errorf("%s loader not implemented yet", loader)
			return http.Dir(""), err
		}
	}

	// fall through to default filesystem loader
	return NewFileSystemStorage(uri)
}

func NewFileSystemStorage(dir string) (http.Dir, error) {

	src, err := os.Stat(dir)
	if err != nil {
		return http.Dir(""), err
	}

	if !src.IsDir() {
		fmt.Errorf("%s is not a directory", dir)
		return http.Dir(""), err
	}

	return http.Dir(dir), nil
}
