package matchers

import "strings"

// HasSuffix returns true when matcher argument ends with the given suffix.
func HasSuffix(value string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.HasSuffix(v, value), nil
	}
}
