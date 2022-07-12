package expect

import "strings"

// ToHavePrefix returns true if the matcher argument starts with the given prefix.
func ToHavePrefix(prefix string) Matcher {
	m := Matcher{}
	m.Name = "HasPrefix"
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.HasPrefix(v.(string), prefix), nil
	}

	return m
}
