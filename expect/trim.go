package expect

import "strings"

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "Trim"
	m.Matches = func(v any, params Args) (bool, error) {
		return matcher.Matches(strings.TrimSpace(v.(string)), params)
	}

	return m
}
