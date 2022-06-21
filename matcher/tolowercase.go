package matcher

import "strings"

func ToLowerCase(m Matcher[string]) Matcher[string] {
	return func(v string, params Params) (bool, error) {
		return m(strings.ToLower(v), params)
	}
}
