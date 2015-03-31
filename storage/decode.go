package storage

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"path/filepath"

	"github.com/disintegration/gift"
)

type ImsImage interface {
	ApplyFilters(*gift.GIFT) error
	Encode() (io.ReadSeeker, error)
}

type GIF struct {
	G *gif.GIF
}

func (j *GIF) ApplyFilters(g *gift.GIFT) error {
	newImages := []*image.Paletted{}
	for _, i := range j.G.Image {
		dst := image.NewPaletted(g.Bounds(i.Bounds()), i.Palette)
		g.Draw(dst, i)
		newImages = append(newImages, dst)
	}

	j.G.Image = newImages
	return nil
}

func (j *GIF) Encode() (io.ReadSeeker, error) {
	bb := new(bytes.Buffer)
	err := gif.EncodeAll(bb, j.G)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}

type JPEG struct {
	I image.Image
}

func (j *JPEG) ApplyFilters(g *gift.GIFT) error {
	dst := image.NewRGBA(g.Bounds(j.I.Bounds()))
	g.Draw(dst, j.I)

	j.I = dst
	return nil
}

func (j *JPEG) Encode() (io.ReadSeeker, error) {
	bb := new(bytes.Buffer)
	err := jpeg.Encode(bb, j.I, &jpeg.Options{Quality: 95})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}

type PNG struct {
	I image.Image
}

func (j *PNG) ApplyFilters(g *gift.GIFT) error {
	dst := image.NewRGBA(g.Bounds(j.I.Bounds()))
	g.Draw(dst, j.I)

	j.I = dst
	return nil
}

func (j *PNG) Encode() (io.ReadSeeker, error) {
	bb := new(bytes.Buffer)
	err := png.Encode(bb, j.I)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}

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
