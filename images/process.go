package images

import (
	"bytes"
	"time"

	"github.com/paddycarey/ims/filters"
	"github.com/paddycarey/ims/storage"
)

func Process(f storage.File, qs string) (storage.File, error) {

	if len(qs) != 0 {

		img, _, err := DecodeImage(f)
		if err != nil {
			return nil, err
		}

		g, pxlBlnd, err := filters.ParseFilters(qs)
		if err != nil {
			return nil, err
		}

		err = img.ApplyFilters(g, pxlBlnd)
		if err != nil {
			return nil, err
		}

		rs, err := img.Encode()
		if err != nil {
			return nil, err
		}

		f = &storage.InMemoryFile{rs.(*bytes.Reader), f.MimeType(), time.Now()}

	}

	return f, nil
}
