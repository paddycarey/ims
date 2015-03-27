package processors

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"

	"github.com/disintegration/gift"
)

type Processor func(*gift.GIFT, string)

// GetProcessor matches strings to processors, useful for mapping URL params to
// actions.
func GetProcessor(name string) (Processor, error) {

	switch name {
	case "crop":
		return Crop, nil
	case "resize":
		return Resize, nil
	}
	return nil, errors.New("Matching processor not found")

}

// ProcessImage loads a file from storage, decodes it as an image, parses the
// query string, processes/filters the image and returns the result wrapped in
// a bytes.NewReader instance.
func ProcessImage(r *http.Request, f http.File) (io.ReadSeeker, error) {

	src, format, err := DecodeImage(f)
	if err != nil {
		log.Println("Unable to decode image")
		log.Printf("%s\n", err)
		return f, nil
	}
	log.Println(format)

	// apply processors to image
	g := gift.New()
	qs := QueryString(r.URL.RawQuery)
	processorArgs, err := qs.Values()
	if err != nil {
		log.Println("Unable to parse query string")
		log.Printf("%s\n", err)
		return f, nil
	}

	for _, procArg := range processorArgs {
		if len(procArg) < 2 {
			continue
		}
		processor, err := GetProcessor(procArg[0])
		if err != nil {
			continue
		}
		processor(g, procArg[1])
	}
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	bb := new(bytes.Buffer)
	err = jpeg.Encode(bb, dst, &jpeg.Options{Quality: 95})
	if err != nil {
		return f, err
	}

	return bytes.NewReader(bb.Bytes()), nil
}
