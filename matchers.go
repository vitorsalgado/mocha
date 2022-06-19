package mocha

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/internal/jsonpath"
)

func Anything[E any]() Matcher[E] {
	return func(v E, params MatcherParams) (bool, error) {
		return true, nil
	}
}

func EqualTo[E any](expected E) Matcher[E] {
	return func(v E, params MatcherParams) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

func Equal(expected any) Matcher[any] {
	return func(v any, params MatcherParams) (bool, error) {
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

func JSONPath[V any](p string, matcher Matcher[V]) Matcher[any] {
	return func(v any, params MatcherParams) (bool, error) {
		value, err := jsonpath.Get(p, v)
		if err != nil || value == nil {
			return false, err
		}

		return matcher(value.(V), params)
	}
}
