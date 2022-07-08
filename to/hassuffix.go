package to

import "strings"

// HaveSuffix returns true when matcher argument ends with the given suffix.
func HaveSuffix(suffix string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "HasSuffix"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.HasSuffix(v, suffix), nil
	}

	return m
}
