package matchers

import "strings"

// HasPrefix returns true if the matcher argument starts with the given prefix.
func HasPrefix(value string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.HasPrefix(v, value), nil
	}
}
