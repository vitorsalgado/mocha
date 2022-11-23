package expect

import (
	"fmt"
	"regexp"
)

// RegExpMatcherTypes defines the acceptable generic types of RegExpMatches.
type RegExpMatcherTypes interface {
	string | regexp.Regexp | *regexp.Regexp
}

type RegExpMatcher[T RegExpMatcherTypes] struct {
	Expression T
}

func (m *RegExpMatcher[T]) Name() string {
	return "MatchRegExp"
}

func (m *RegExpMatcher[T]) Match(v any) (bool, error) {
	var err error
	var result bool

	switch e := any(m.Expression).(type) {
	case string:
		return regexp.Match(e, []byte(fmt.Sprintf("%v", v)))
	case regexp.Regexp:
		result = e.Match([]byte(fmt.Sprintf("%v", v)))
	case *regexp.Regexp:
		result = e.Match([]byte(fmt.Sprintf("%v", v)))
	}

	return result, err
}

func (m *RegExpMatcher[T]) DescribeFailure(_ any) string {
	return "given regular expression dit not match"
}

func (m *RegExpMatcher[T]) OnMockServed() error {
	return nil
}

// ToMatchExpr returns true then the given regular expression matches matcher argument.
// ToMatchExpr accepts a string or a regexp.Regexp.
func ToMatchExpr[T RegExpMatcherTypes](re T) Matcher {
	return &RegExpMatcher[T]{Expression: re}
}
