package expect

import "strings"

// ToHaveSuffix returns true when matcher argument ends with the given suffix.
func ToHaveSuffix(suffix string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "HasSuffix"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.HasSuffix(v, suffix), nil
	}

	return m
}
