package to

import "strings"

// HavePrefix returns true if the matcher argument starts with the given prefix.
func HavePrefix(prefix string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "HasPrefix"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.HasPrefix(v, prefix), nil
	}

	return m
}
