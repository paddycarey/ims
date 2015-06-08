package images

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/gift"
)

type ImsImage interface {
	ApplyFilters(*gift.GIFT, bool) error
	Encode() (io.ReadSeeker, error)
}

type GIF struct {
	G *gif.GIF
}

func (j *GIF) ApplyFilters(g *gift.GIFT, pb bool) error {
	newImages := []*image.Paletted{}
	firstFrame := j.G.Image[0]

	if !pb {
		for i := range j.G.Image {
			// tmp image is used here to keep the same dimensions for each frame
			tmp := image.NewNRGBA(firstFrame.Bounds())
			gift.New().DrawAt(tmp, j.G.Image[i], j.G.Image[i].Bounds().Min, gift.CopyOperator)
			dst := image.NewPaletted(g.Bounds(tmp.Bounds()), j.G.Image[i].Palette)
			g.Draw(dst, tmp)
			newImages = append(newImages, dst)
		}
	} else {
		tmp := image.NewNRGBA(firstFrame.Bounds())
		for i := range j.G.Image {
			// draw current frame over previous:
			gift.New().DrawAt(tmp, j.G.Image[i], j.G.Image[i].Bounds().Min, gift.OverOperator)
			dst := image.NewPaletted(g.Bounds(tmp.Bounds()), j.G.Image[i].Palette)
			g.Draw(dst, tmp)
			newImages = append(newImages, dst)
		}
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

func (j *JPEG) ApplyFilters(g *gift.GIFT, pb bool) error {
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

func (j *PNG) ApplyFilters(g *gift.GIFT, pb bool) error {
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
