package processors

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"path/filepath"
)

// DecodeImage reads a file from storage and attempts to decode it as an
// image.Image. We first attempt to autodetect the encoding and decode
// automatically, if this fails then we attempt to decode manually using a
// decoder based on the file extension.
func DecodeImage(f http.File) (image.Image, string, error) {

	// attempt to decode the image automatically
	src, format, err := image.Decode(f)
	if err == nil {
		return src, format, err
	}

	// if we can't decode the image automatically then we'll try to manually
	// decode based on the file's extension.
	fi, err := f.Stat()
	if err != nil {
		return src, "", err
	}
	ext := filepath.Ext(fi.Name())
	switch ext {
	case ".gif":
		src, err = gif.Decode(f)
		return src, "gif", err
	case ".jpeg", ".jpg":
		src, err = jpeg.Decode(f)
		return src, "jpeg", err
	case ".png":
		src, err = png.Decode(f)
		return src, "png", err
	}

	err = errors.New("Unable to decode image")
	return src, "", err
}
