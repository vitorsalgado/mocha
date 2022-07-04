package matchers

import "strings"

// HasPrefix returns true if the matcher argument starts with the given prefix.
func HasPrefix(prefix string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "HasPrefix"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.HasPrefix(v, prefix), nil
	}

	return m
}
