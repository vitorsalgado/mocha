package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type splitMatcher struct {
	Separator string
	Matcher   Matcher
}

func (m *splitMatcher) Name() string {
	return fmt.Sprintf("Split(%s)", m.Matcher.Name())
}

func (m *splitMatcher) Match(v any) (*Result, error) {
	txt, ok := v.(string)
	if !ok {
		return &Result{}, fmt.Errorf("type %s is not supported. only string is acceptable", reflect.TypeOf(v).Name())
	}

	result, err := m.Matcher.Match(strings.Split(txt, m.Separator))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), printExpected(txt)),
				result.DescribeFailure(),
			)
		},
	}, nil
}

func (m *splitMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

func Split(separator string, matcher Matcher) Matcher {
	return &splitMatcher{Separator: separator, Matcher: matcher}
}
