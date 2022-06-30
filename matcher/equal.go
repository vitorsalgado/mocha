package matcher

import (
	"reflect"
	"strings"
)

func EqualTo[E any](expected E) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

func EqualAny(expected any) Matcher[any] {
	return func(v any, params Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

func EqualFold(expected string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.EqualFold(expected, v), nil
	}
}
