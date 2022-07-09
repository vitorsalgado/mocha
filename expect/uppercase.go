package expect

import "strings"

// UpperCase upper case matcher string argument before submitting it to provided matcher.
func UpperCase(matcher Matcher[string]) Matcher[string] {
	m := Matcher[string]{}
	m.Name = "ToUpperCase"
	m.Matches =
		func(v string, params Args) (bool, error) {
			return matcher.Matches(strings.ToUpper(v), params)
		}

	return m
}
