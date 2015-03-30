package filters

import (
	"image"
	"strconv"

	"github.com/disintegration/gift"
)

type Filter func(*gift.GIFT, []string)

// Crop will crop an image at the specified coords (x0, y0, x1, y1)
func Crop(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 4 {
		return
	}

	cropInt := []int{}
	for _, s := range qsArgs {
		i, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		if i < 0 {
			return
		}
		cropInt = append(cropInt, i)
	}

	g.Add(gift.Crop(image.Rect(cropInt[0], cropInt[1], cropInt[2], cropInt[3])))
}

// Resize resizes the given image to the specified size with the configured
// resampling algorithm (defaults to lanczos).
func Resize(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) < 2 || len(qsArgs) > 3 {
		return
	}

	w, err := strconv.Atoi(qsArgs[0])
	if err != nil {
		return
	}

	h, err := strconv.Atoi(qsArgs[1])
	if err != nil {
		return
	}

	if w == 0 && h == 0 {
		return
	}

	sampler := gift.LanczosResampling
	if len(qsArgs) > 2 {
		switch qsArgs[2] {
		case "box":
			sampler = gift.BoxResampling
		case "cubic":
			sampler = gift.CubicResampling
		case "lanczos":
			sampler = gift.LanczosResampling
		case "linear":
			sampler = gift.LinearResampling
		case "nearestneighbour":
			sampler = gift.NearestNeighborResampling
		}
	}

	g.Add(gift.Resize(w, h, sampler))
}
