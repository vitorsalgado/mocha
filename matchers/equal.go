package matchers

import "reflect"

// EqualTo returns true if matcher value is equal to the given parameter value.
func EqualTo[E any](expected E) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}

// EqualAny is equivalent to the EqualTo function, but it is not generic.
func EqualAny(expected any) Matcher[any] {
	return func(v any, params Args) (bool, error) {
		return reflect.DeepEqual(expected, v), nil
	}
}
