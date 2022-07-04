package matchers

import "strings"

// ToUpperCase upper case matcher string argument before submitting it to provided matcher.
func ToUpperCase(m Matcher[string]) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return m(strings.ToUpper(v), params)
	}
}
