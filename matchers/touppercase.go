package matchers

import "strings"

// ToUpperCase upper case matcher string argument before submitting it to provided matcher.
func ToUpperCase(matcher Matcher[string]) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "ToUpperCase"
	m.Matches =
		func(v string, params Args) (bool, error) {
			return matcher.Matches(strings.ToUpper(v), params)
		}

	return m
}
