package mocha

import (
	"net/url"
	"reflect"
	"strings"
)

func Anything[E comparable]() Matcher[E] {
	return func(v E, params MatcherParams) (bool, error) {
		return true, nil
	}
}

func Equal[E any](expected E) Matcher[E] {
	return func(v E, params MatcherParams) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

func EqualFold(expected string) Matcher[string] {
	return func(v string, params MatcherParams) (bool, error) {
		return strings.EqualFold(expected, v), nil
	}
}

func URLPath(expected string) Matcher[url.URL] {
	return func(v url.URL, params MatcherParams) (bool, error) {
		return strings.EqualFold(expected, v.Path), nil
	}
}
