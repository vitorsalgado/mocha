package expect

import "strings"

// UpperCase upper case matcher string argument before submitting it to provided matcher.
func UpperCase(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "UpperCase"
	m.Matches =
		func(v any, params Args) (bool, error) {
			return matcher.Matches(strings.ToUpper(v.(string)), params)
		}

	return m
}
