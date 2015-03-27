package processors

import (
	"image"
	"strconv"
	"strings"

	"github.com/disintegration/gift"
)

// Crop will crop an image at the specified coords (x0, y0, x1, y1)
func Crop(g *gift.GIFT, qsArgs string) {

	cropStr := strings.Split(qsArgs, ",")
	if len(cropStr) != 4 {
		return
	}

	cropInt := []int{}
	for _, s := range cropStr {
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
