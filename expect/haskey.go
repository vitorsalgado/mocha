package expect

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

// ToHaveKey returns true if the JSON key in the given path is present.
// Example:
//
//	JSON: { "name": "test" }
//	ToHaveKey("name") will return true
//	ToHaveKey("address.street") will return false.
func ToHaveKey(path string) Matcher {
	m := Matcher{}
	m.Name = "HasKey"
	m.DescribeMismatch = func(p string, v any) string {
		return fmt.Sprintf("json doest not have a key on path: %s", path)
	}
	m.Matches = func(v any, args Args) (bool, error) {
		value, err := jsonx.Reach(path, v)
		if err != nil || value == nil {
			return false, err
		}

		return true, nil
	}

	return m
}
