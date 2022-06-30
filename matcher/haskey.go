package matcher

import "github.com/vitorsalgado/mocha/internal/jsonpath"

func HasKey[V any](path string) Matcher[any] {
	return func(v any, params Args) (bool, error) {
		value, err := jsonpath.Get(path, v)
		if err != nil || value == nil {
			return false, err
		}

		return true, nil
	}
}
