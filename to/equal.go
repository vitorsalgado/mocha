package to

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/inspect"
)

// Equal returns true if matcher value is equal to the given parameter value.
func Equal[V any](expected V) Matcher[V] {
	matcher := Matcher[V]{}
	matcher.Name = "equalTo"
	matcher.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(inspect.ToString(expected))),
			fmt.Sprintf("got: %s", colorize.Red(inspect.ToString(v))),
		)
	}
	matcher.Matches = func(v V, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return matcher
}

// EqualAny is equivalent to the EqualTo function, but it is not generic.
func EqualAny(expected any) Matcher[any] {
	m := Matcher[any]{}
	m.Name = "EqualAny"
	m.Matches = func(v any, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return m
}
