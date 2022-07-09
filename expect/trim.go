package expect

import "strings"

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher[string]) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "Trim"
	m.Matches = func(v string, params Args) (bool, error) {
		return matcher.Matches(strings.TrimSpace(v), params)
	}

	return m
}
