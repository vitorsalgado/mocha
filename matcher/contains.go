package matcher

import "strings"

func Contains(value string) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return strings.Contains(v, value), nil
	}
}
