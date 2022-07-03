package matcher

import "strings"

// EqualFold returns true if expected value is equal to matcher value, ignoring case.
// EqualFold uses strings.EqualFold function.
func EqualFold(expected string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.EqualFold(expected, v), nil
	}
}
