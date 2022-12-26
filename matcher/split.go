package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type splitMatcher struct {
	separator string
	matcher   Matcher
}

func (m *splitMatcher) Name() string {
	return fmt.Sprintf("Split(%s)", m.matcher.Name())
}

func (m *splitMatcher) Match(v any) (*Result, error) {
	txt, ok := v.(string)
	if !ok {
		return &Result{}, fmt.Errorf("type %s is not supported. only string is acceptable", reflect.TypeOf(v).Name())
	}

	result, err := m.matcher.Match(strings.Split(txt, m.separator))
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		Pass: result.Pass,
		Message: func() string {
			return fmt.Sprintf("%s %s",
				hint(m.Name(), printExpected(txt)),
				result.Message(),
			)
		},
	}, nil
}

func (m *splitMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

func Split(separator string, matcher Matcher) Matcher {
	return &splitMatcher{separator: separator, matcher: matcher}
}
