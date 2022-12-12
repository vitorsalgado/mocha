package matcher

import (
	"fmt"
	"reflect"
	"regexp"
)

type regExpMatcher struct {
	expression any
}

func (m *regExpMatcher) Name() string {
	return "MatchRegExp"
}

func (m *regExpMatcher) Match(v any) (*Result, error) {
	txt := fmt.Sprintf("%v", v)

	msg := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.expression)),
			_separator,
			txt)
	}

	switch e := m.expression.(type) {
	case string:
		match, err := regexp.Match(e, []byte(txt))
		return &Result{OK: match, DescribeFailure: msg}, err
	case regexp.Regexp:
		return &Result{OK: e.Match([]byte(txt)), DescribeFailure: msg}, nil
	case *regexp.Regexp:
		return &Result{OK: e.Match([]byte(txt)), DescribeFailure: msg}, nil
	default:
		return mismatch(nil), fmt.Errorf("regular expression matcher does not accept the expression of type %s",
			reflect.TypeOf(v).Name())
	}
}

func (m *regExpMatcher) OnMockServed() error {
	return nil
}

func (m *regExpMatcher) Spec() any {
	return []any{_mRegex, m.expression}
}

// Matches returns true then the given regular expression matches matcher argument.
// It accepts a string or a regexp.Regexp.
func Matches(expression any) Matcher {
	return &regExpMatcher{expression: expression}
}
