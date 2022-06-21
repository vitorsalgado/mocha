package matcher

import "strings"

func Trim(ms Matcher[string]) Matcher[string] {
	return func(v string, params Params) (bool, error) {
		return ms(strings.TrimSpace(v), params)
	}
}
