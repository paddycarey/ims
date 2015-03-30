package filters

import (
	"image"

	"github.com/disintegration/gift"
)

func ApplyFilters(g *gift.GIFT, i image.Image) image.Image {

	// apply processors to image
	dst := image.NewRGBA(g.Bounds(i.Bounds()))
	g.Draw(dst, i)

	return dst
}
