package expect

import (
	"github.com/vitorsalgado/mocha/x/jsonx"
)

// ToHaveKey returns true if the JSON key in the given path is present.
// Example:
//	JSON: { "name": "test" }
//	ToHaveKey("name") will return true
//	ToHaveKey("address.street") will return false.
func ToHaveKey[V any](path string) Matcher[any] {
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
