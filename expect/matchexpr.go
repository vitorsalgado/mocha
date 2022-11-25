package expect

import (
	"fmt"
	"reflect"
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

func (m *RegExpMatcher[T]) Match(v any) (Result, error) {
	txt := fmt.Sprintf("%v", v)

	msg := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.Expression)),
			_separator,
			txt)
	}

	switch e := any(m.Expression).(type) {
	case string:
		match, err := regexp.Match(e, []byte(txt))
		return Result{OK: match, DescribeFailure: msg}, err
	case regexp.Regexp:
		return Result{OK: e.Match([]byte(txt)), DescribeFailure: msg}, nil
	case *regexp.Regexp:
		return Result{OK: e.Match([]byte(txt)), DescribeFailure: msg}, nil
	default:
		return mismatch(nil),
			fmt.Errorf("regular expression matcher does not accept the expression of type %s",
				reflect.TypeOf(v).Name())
	}
}

func (m *RegExpMatcher[T]) OnMockServed() error {
	return nil
}

// ToMatchExpr returns true then the given regular expression matches matcher argument.
// ToMatchExpr accepts a string or a regexp.Regexp.
func ToMatchExpr[T RegExpMatcherTypes](re T) Matcher {
	return &RegExpMatcher[T]{Expression: re}
}
