package filters

import (
	"fmt"
	"image"
	"strconv"

	"github.com/disintegration/gift"
)

type Filter func(*gift.GIFT, []string)

// GetFilter returns an appropriate Filter function for a given string.
func GetFilter(name string) (Filter, error) {

	switch name {
	case "brightness":
		return Brightness, nil
	case "contrast":
		return Contrast, nil
	case "crop":
		return Crop, nil
	case "fliphorizontal":
		return FlipHorizontal, nil
	case "flipvertical":
		return FlipVertical, nil
	case "hue":
		return Hue, nil
	case "resize":
		return Resize, nil
	case "rotate":
		return Rotate, nil
	case "saturation":
		return Saturation, nil
	case "transpose":
		return Transpose, nil
	case "transverse":
		return Transverse, nil
	}

	return nil, fmt.Errorf("Matching filter not found: %s", name)
}

// Brightness creates a filter that changes the brightness of an image. The
// first (and only) parameter (percentage) must be in range (-100, 100). The
// percentage = 0 gives the original image. The percentage = -100 gives solid
// black image. The percentage = 100 gives solid white image.
func Brightness(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 1 {
		return
	}

	i64, err := strconv.ParseFloat(qsArgs[0], 32)
	if err != nil {
		return
	}
	i32 := float32(i64)
	if i32 < -100 || i32 > 100 {
		return
	}

	g.Add(gift.Brightness(i32))

}

// Contrast creates a filter that changes the contrast of an image. The first
// (and only) parameter (percentage) must be in range (-100, 100). The
// percentage = 0 gives the original image. The percentage = -100 gives solid
// grey image. The percentage = 100 gives an overcontrasted image.
func Contrast(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 1 {
		return
	}

	i64, err := strconv.ParseFloat(qsArgs[0], 32)
	if err != nil {
		return
	}
	i32 := float32(i64)
	if i32 < -100 || i32 > 100 {
		return
	}

	g.Add(gift.Contrast(i32))

}

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

// FlipHorizontal creates a filter that flips an image horizontally.
func FlipHorizontal(g *gift.GIFT, qsArgs []string) {
	g.Add(gift.FlipHorizontal())
}

// FlipVertical creates a filter that flips an image vertically.
func FlipVertical(g *gift.GIFT, qsArgs []string) {
	g.Add(gift.FlipVertical())
}

// Hue creates a filter that rotates the hue of an image. The shift parameter
// is the hue angle shift, typically in range (-180, 180). The shift = 0 gives
// the original image.
func Hue(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 1 {
		return
	}

	i64, err := strconv.ParseFloat(qsArgs[0], 32)
	if err != nil {
		return
	}
	i32 := float32(i64)
	if i32 < -180 || i32 > 180 {
		return
	}

	g.Add(gift.Hue(i32))

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

// Rotate creates a filter that rotates an image by the given angle
// counter-clockwise. The angle parameter is the rotation angle in degrees. The
// only allowed values for this fuction are 90, 180 and 270
func Rotate(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 1 {
		return
	}

	switch qsArgs[0] {
	case "90":
		g.Add(gift.Rotate90())
	case "180":
		g.Add(gift.Rotate180())
	case "270":
		g.Add(gift.Rotate270())
	}

}

// Saturation creates a filter that changes the saturation of an image. The
// percentage parameter must be in range (-100, 500). The percentage = 0 gives
// the original image.
func Saturation(g *gift.GIFT, qsArgs []string) {

	if len(qsArgs) != 1 {
		return
	}

	i64, err := strconv.ParseFloat(qsArgs[0], 32)
	if err != nil {
		return
	}
	i32 := float32(i64)
	if i32 < -100 || i32 > 500 {
		return
	}

	g.Add(gift.Saturation(i32))

}

// Transpose creates a filter that flips an image horizontally and rotates 90
// degrees counter-clockwise.
func Transpose(g *gift.GIFT, qsArgs []string) {
	g.Add(gift.Transpose())
}

// Transverse creates a filter that flips an image vertically and rotates 90
// degrees counter-clockwise.
func Transverse(g *gift.GIFT, qsArgs []string) {
	g.Add(gift.Transverse())
}
