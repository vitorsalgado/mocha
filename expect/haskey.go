package expect

import (
	"github.com/vitorsalgado/mocha/util/jsonx"
)

// ToHaveProperty returns true if the JSON key in the given path is present.
// Example:
//	JSON: { "name": "test" }
//	ToHaveProperty("name") will return true
//	ToHaveProperty("address.street") will return false.
func ToHaveProperty[V any](path string) Matcher[any] {
	m := Matcher[any]{}
	m.Name = "HasKey"
	m.Matches = func(v any, args Args) (bool, error) {
		value, err := jsonx.Reach(path, v)
		if err != nil || value == nil {
			return false, err
		}

		return true, nil
	}

	return m
}
