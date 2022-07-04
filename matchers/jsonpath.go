package matchers

import "github.com/vitorsalgado/mocha/internal/jsonpath"

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath[V any](p string, matcher Matcher[V]) Matcher[any] {
	m := Matcher[any]{}
	m.Name = "JSONPath"
	m.Matches = func(v any, args Args) (bool, error) {
		value, err := jsonpath.Get(p, v)
		if err != nil || value == nil {
			return false, err
		}

		return matcher.Matches(value.(V), args)
	}

	return m
}
