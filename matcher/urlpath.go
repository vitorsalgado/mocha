package matcher

import (
	"net/url"
	"strings"
)

func URLPath(expected string) Matcher[url.URL] {
	return func(v url.URL, params Args) (bool, error) {
		return strings.EqualFold(expected, v.Path), nil
	}
}
