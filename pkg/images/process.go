package images

import (
	"bytes"
	"net/url"
	"time"

	"github.com/Sirupsen/logrus"

	"github.com/paddycarey/ims/pkg/filters"
	"github.com/paddycarey/ims/pkg/storage"
)

func Process(f storage.File, u *url.URL, noOpts bool) (storage.File, error) {

	if len(u.RawQuery) != 0 {

		logrus.WithField("url", u.RequestURI()).Debug("Processing image")

		img, _, err := DecodeImage(f)
		if err != nil {
			return nil, err
		}

		g, pxlBlnd, err := filters.ParseFilters(u.RawQuery)
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

	// optimise unless the user has explicitly disabled it
	if !noOpts {
		logrus.WithField("url", u.RequestURI()).Debug("Optimizing image")
		var err error
		f, err = Optimize(f)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}
