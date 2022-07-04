package matchers

import (
	"net/url"
	"strings"
)

// URLPath returns true if request URL path is equal to the expected path, ignoring case.
func URLPath(expected string) Matcher[url.URL] {
	return func(v url.URL, params Args) (bool, error) {
		return strings.EqualFold(expected, v.Path), nil
	}
}
