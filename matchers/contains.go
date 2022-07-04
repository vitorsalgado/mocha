package matchers

import "strings"

// Contains returns true when the expected value is contained in the matcher argument.
func Contains(value string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.Contains(v, value), nil
	}
}
