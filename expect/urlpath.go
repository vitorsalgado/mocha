package expect

import (
	"net/url"
	"strings"
)

// URLPath returns true if request URL path is equal to the expected path, ignoring case.
func URLPath(expected string) Matcher {
	m := Matcher{}
	m.Name = "URLPath"
	m.Matches = func(v any, params Args) (bool, error) {
		switch e := v.(type) {
		case *url.URL:
			return strings.EqualFold(expected, e.Path), nil
		case url.URL:
			return strings.EqualFold(expected, e.Path), nil
		case string:
			u, err := url.Parse(e)
			if err != nil {
				return false, err
			}

			return strings.EqualFold(expected, u.Path), nil

		default:
			panic("URLPath matcher only accepts the types: *url.URL | url.URL | string")
		}
	}

	return m
}
