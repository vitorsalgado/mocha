package matcher

import (
	"fmt"
	"regexp"
)

type RegExpMatcherTypes interface {
	string | regexp.Regexp | *regexp.Regexp
}

func RegExpMatches[V any, T RegExpMatcherTypes](re T) Matcher[V] {
	return func(v V, params Args) (bool, error) {
		switch e := any(re).(type) {
		case string:
			return regexp.Match(e, []byte(fmt.Sprintf("%v", v)))
		case regexp.Regexp:
			return e.Match([]byte(fmt.Sprintf("%v", v))), nil
		case *regexp.Regexp:
			return e.Match([]byte(fmt.Sprintf("%v", v))), nil
		}

		return false, fmt.Errorf("unable to apply regexp expression for value %v", v)
	}
}
