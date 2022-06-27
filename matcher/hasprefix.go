package matcher

import "strings"

func HasPrefix(value string) Matcher[string] {
	return func(v string, params Params) (bool, error) {
		return strings.HasPrefix(v, value), nil
	}
}
