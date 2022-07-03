package matcher

import "strings"

// ToLowerCase lower case matcher string argument before submitting it to provided matcher.
func ToLowerCase(m Matcher[string]) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return m(strings.ToLower(v), params)
	}
}
