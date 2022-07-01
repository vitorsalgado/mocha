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
		var err error
		var result bool

		switch e := any(re).(type) {
		case string:
			return regexp.Match(e, []byte(fmt.Sprintf("%v", v)))
		case regexp.Regexp:
			result = e.Match([]byte(fmt.Sprintf("%v", v)))
		case *regexp.Regexp:
			result = e.Match([]byte(fmt.Sprintf("%v", v)))
		}

		return result, err
	}
}
