package expect

import "strings"

// ToHaveSuffix returns true when matcher argument ends with the given suffix.
func ToHaveSuffix(suffix string) Matcher {
	m := Matcher{}
	m.Name = "HasSuffix"
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.HasSuffix(v.(string), suffix), nil
	}

	return m
}
