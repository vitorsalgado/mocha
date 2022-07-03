package matcher

import "strings"

// Trim trims matcher string arguemtn before submitting it to the provided matcher.
func Trim(ms Matcher[string]) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return ms(strings.TrimSpace(v), params)
	}
}
