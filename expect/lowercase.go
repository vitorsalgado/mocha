package expect

import "strings"

// LowerCase lower case matcher string argument before submitting it to provided matcher.
func LowerCase(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "LowerCase"
	m.Matches = func(v any, params Args) (bool, error) {
		return matcher.Matches(strings.ToLower(v.(string)), params)
	}

	return m
}
