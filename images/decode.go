package images

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
func DecodeImage(f http.File) (ImsImage, string, error) {

	// attempt to decode the image automatically
	src, format, err := image.Decode(f)
	if err == nil {
		switch format {
		case "jpeg":
			return &JPEG{src}, "jpeg", nil
		case "png":
			return &PNG{src}, "png", nil
		}
	}

	// if we can't decode the image automatically then we'll try to manually
	// decode based on the file's extension.
	fi, err := f.Stat()
	if err != nil {
		return nil, "", err
	}
	// seek back to beginning of file so we can decode again if we need to
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, "", err
	}
	ext := filepath.Ext(fi.Name())
	switch ext {
	case ".gif":
		gifSrc, err := gif.DecodeAll(f)
		if err != nil {
			return nil, "", err
		}
		return &GIF{gifSrc}, "gif", err
	case ".jpeg", ".jpg":
		src, err = jpeg.Decode(f)
		if err != nil {
			return nil, "", err
		}
		return &JPEG{src}, "jpeg", err
	case ".png":
		src, err = png.Decode(f)
		if err != nil {
			return nil, "", err
		}
		return &PNG{src}, "png", err
	}

	err = errors.New("Unable to decode image")
	return nil, "", err
}
