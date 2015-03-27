package processors

import (
	"net/url"
	"strings"
)

// QueryString is a convenience wrapper designed to be wrapped around a
// url.RawQuery, providing ordered access to all key/value pairs.
type QueryString string

// Values returns a slice of []string, each []string will have a length of 2
// and will represent a single key/value pair parsed from the query string.
func (qs *QueryString) Values() ([][]string, error) {

	var err error
	query := string(*qs)
	values := [][]string{}
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
		values = append(values, []string{key, value})
	}
	return values, err
}
