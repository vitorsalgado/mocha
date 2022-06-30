package matcher

import "strings"

func ToLowerCase(m Matcher[string]) Matcher[string] {
	return func(v string, params Args) (bool, error) {
		return m(strings.ToLower(v), params)
	}
}
