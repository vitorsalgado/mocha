package matcher

import (
	"strings"
)

func Contains[E any](value E) Matcher[E] {
	return func(v E, params Params) (bool, error) {
		switch e := any(value).(type) {
		case string:
			return strings.Contains(any(v).(string), e), nil
		}

		return false, nil
	}
}
