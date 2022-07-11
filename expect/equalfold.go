package expect

import "strings"

// ToEqualFold returns true if expected value is equal to matcher value, ignoring case.
// ToEqualFold uses strings.EqualFold function.
func ToEqualFold(expected string) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "EqualFold"
	m.Matches = func(v string, args Args) (bool, error) {
		return strings.EqualFold(expected, v), nil
	}

	return m
}
