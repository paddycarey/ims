package filters

import (
	"net/url"
	"strings"

	"github.com/disintegration/gift"
)

type FilterArgs struct {
	Filter string
	Args   []string
}

func NewFilterArgs(filter, args string) *FilterArgs {
	return &FilterArgs{
		Filter: filter,
		Args:   strings.Split(args, ","),
	}
}

// QueryString is a convenience wrapper designed to be wrapped around a
// url.RawQuery, providing ordered access to all key/value pairs.
type QueryString string

// Values returns a slice of *FilterArgs, each []string will have a length of 2
// and will represent a single key/value pair parsed from the query string.
func (qs *QueryString) Values() ([]*FilterArgs, error) {

	var err error
	query := string(*qs)
	values := []*FilterArgs{}
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		values = append(values, NewFilterArgs(key, value))
	}
	return values, err
}

// ProcessImage loads a file from storage, decodes it as an image, parses the
// query string, processes/filters the image and returns the result wrapped in
// a bytes.NewReader instance.
func ParseFilters(q string) (*gift.GIFT, bool, error) {

	// apply processors to image
	qs := QueryString(q)
	filterArgs, err := qs.Values()
	if err != nil {
		return nil, false, err
	}

	g := gift.New()
	pxlBlnd := false
	for _, filterArg := range filterArgs {
		filter, err := GetFilter(filterArg.Filter)
		if err != nil {
			return nil, false, err
		}
		fPxlBlnd := filter(g, filterArg.Args)
		if fPxlBlnd {
			pxlBlnd = true
		}
	}

	return g, pxlBlnd, nil
}
