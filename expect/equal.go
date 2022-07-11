package expect

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/internal/colorize"
	"github.com/vitorsalgado/mocha/internal/misc"
)

// ToEqual returns true if matcher value is equal to the given parameter value.
func ToEqual[V any](expected V) Matcher[V] {
	matcher := Matcher[V]{}
	matcher.Name = "EqualTo"
	matcher.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("%s\n%s",
			fmt.Sprintf("expected: %v", colorize.Green(misc.ToString(expected))),
			fmt.Sprintf("got: %s", colorize.Red(misc.ToString(v))),
		)
	}
	matcher.Matches = func(v V, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return matcher
}

// ToEqualAny is equivalent to the EqualTo function, but it is not generic.
func ToEqualAny(expected any) Matcher[any] {
	m := Matcher[any]{}
	m.Name = "EqualAny"
	m.Matches = func(v any, args Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}

	return m
}
