package to

import "strings"

// LowerCase lower case matcher string argument before submitting it to provided matcher.
func LowerCase(matcher Matcher[string]) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "ToLowerCase"
	m.Matches = func(v string, params Args) (bool, error) {
		return matcher.Matches(strings.ToLower(v), params)
	}

	return m
}
