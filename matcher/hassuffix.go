package matcher

import "strings"

func HasSuffix(value string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.HasSuffix(v, value), nil
	}
}
