package filters

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

type Encoder func(image.Image) (io.ReadSeeker, error)

func GetEncoder(format string) (Encoder, error) {

	switch format {
	case "gif":
		return EncodeGIF, nil
	case "jpeg":
		return EncodeJPEG, nil
	case "png":
		return EncodePNG, nil
	}

	return nil, fmt.Errorf("Matching encoder not found: %s", format)
}

func EncodeGIF(i image.Image) (io.ReadSeeker, error) {

	bb := new(bytes.Buffer)
	err := gif.Encode(bb, i, &gif.Options{})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}

func EncodeJPEG(i image.Image) (io.ReadSeeker, error) {

	bb := new(bytes.Buffer)
	err := jpeg.Encode(bb, i, &jpeg.Options{Quality: 95})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}

func EncodePNG(i image.Image) (io.ReadSeeker, error) {

	bb := new(bytes.Buffer)
	err := png.Encode(bb, i)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}
