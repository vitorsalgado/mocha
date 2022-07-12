package expect

import (
	"fmt"
	"regexp"
)

// RegExpMatcherTypes defines the acceptable generic types of RegExpMatches.
type RegExpMatcherTypes interface {
	string | regexp.Regexp | *regexp.Regexp
}

// ToMatchExpr returns true then the given regular expression matches matcher argument.
// ToMatchExpr accepts a string or a regexp.Regexp.
func ToMatchExpr[T RegExpMatcherTypes](re T) Matcher {
	m := Matcher{}
	m.Name = "MatchExpr"
	m.Matches = func(v any, params Args) (bool, error) {
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

	return m
}
