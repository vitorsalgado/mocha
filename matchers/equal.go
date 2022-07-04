package matchers

import "reflect"

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo[V any](expected V) Matcher[V] {
	matcher := Matcher[V]{}
	matcher.Name = "equalTo"
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
