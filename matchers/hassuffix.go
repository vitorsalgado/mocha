package matchers

import "strings"

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(suffix string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "HasSuffix"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.HasSuffix(v, suffix), nil
	}

	return m
}
