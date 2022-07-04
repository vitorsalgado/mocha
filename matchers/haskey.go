package matchers

import "github.com/vitorsalgado/mocha/internal/jsonpath"

// HasKey returns true if the JSON key in the given path is present.
// Example:
//	JSON: { "name": "test" }
//	HasKey("name") will return true
//	HasKey("address.street") will return false.
func HasKey[V any](path string) Matcher[any] {
	m := Matcher[any]{}
	m.Name = "HasKey"
	m.Matches = func(v any, args Args) (bool, error) {
		value, err := jsonpath.Get(path, v)
		if err != nil || value == nil {
			return false, err
		}

		return true, nil
	}

	return m
}
