package to

import "github.com/vitorsalgado/mocha/internal/jsonpath"

// HaveProperty returns true if the JSON key in the given path is present.
// Example:
//	JSON: { "name": "test" }
//	HaveProperty("name") will return true
//	HaveProperty("address.street") will return false.
func HaveProperty[V any](path string) Matcher[any] {
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
