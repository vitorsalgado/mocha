package matcher

import "strings"

func HasSuffix(value string) Matcher[string] {
	return func(v string, params Params) (bool, error) {
		return strings.HasSuffix(v, value), nil
	}
}
