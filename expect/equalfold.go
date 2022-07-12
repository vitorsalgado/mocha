package expect

import "strings"

// ToEqualFold returns true if expected value is equal to matcher value, ignoring case.
// ToEqualFold uses strings.EqualFold function.
func ToEqualFold(expected string) Matcher {
	m := Matcher{}
	m.Name = "EqualFold"
	m.Matches = func(v any, args Args) (bool, error) {
		return strings.EqualFold(expected, v.(string)), nil
	}

	return m
}
