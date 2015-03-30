package storage

import (
	"fmt"
	"net/http"
	"os"
)

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
