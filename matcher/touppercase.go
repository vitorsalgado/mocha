package matcher

import "strings"

func ToUpperCase(m Matcher[string]) Matcher[string] {
	return func(v string, params Params) (bool, error) {
		return m(strings.ToUpper(v), params)
	}
}
