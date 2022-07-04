package matchers

import "strings"

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(ms Matcher[string]) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return ms(strings.TrimSpace(v), params)
	}
}
