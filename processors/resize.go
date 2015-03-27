package processors

import (
	"strconv"
	"strings"

	"github.com/disintegration/gift"
)

// Resize resizes the given image to the specified size with the configured
// resampling algorithm (defaults to lanczos).
func Resize(g *gift.GIFT, qsArgs string) {

	resizeStr := strings.Split(qsArgs, ",")
	if len(resizeStr) < 2 || len(resizeStr) > 3 {
		return
	}

	w, err := strconv.Atoi(resizeStr[0])
	if err != nil {
		return
	}

	h, err := strconv.Atoi(resizeStr[1])
	if err != nil {
		return
	}

	if w == 0 && h == 0 {
		return
	}

	sampler := gift.LanczosResampling
	if len(resizeStr) > 2 {
		switch resizeStr[2] {
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
