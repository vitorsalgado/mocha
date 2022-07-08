package to

import (
	"net/url"
	"strings"
)

// HaveURLPath returns true if request URL path is equal to the expected path, ignoring case.
func HaveURLPath(expected string) Matcher[url.URL] {
	m := Matcher[url.URL]{}
	m.Name = "URLPath"
	m.Matches = func(v url.URL, params Args) (bool, error) {
		return strings.EqualFold(expected, v.Path), nil
	}

	return m
}
