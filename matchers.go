package mocha

import (
	"net/url"
	"reflect"
	"strings"
)

func Anything[E comparable]() Matcher[E] {
	return func(v E, ctx MatcherContext) (bool, error) {
		return true, nil
	}
}

func Equal[E any](expected E) Matcher[E] {
	return func(v E, ctx MatcherContext) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

func EqualFold(expected string) Matcher[string] {
	return func(v string, ctx MatcherContext) (bool, error) {
		return strings.EqualFold(expected, v), nil
	}
}

func URLPath(expected string) Matcher[url.URL] {
	return func(v url.URL, ctx MatcherContext) (bool, error) {
		return strings.EqualFold(expected, v.Path), nil
	}
}
