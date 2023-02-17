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
		return nil, fmt.Errorf("type %s is not supported. accepted types: string", reflect.TypeOf(v))
	}

	result, err := m.matcher.Match(strings.Split(txt, m.separator))
	if err != nil {
		return nil, err
	}

	if result.Pass {
		return &Result{Pass: true}, err
	}

	return &Result{
		Ext:     []string{txt},
		Message: result.Message,
	}, nil
}

func (m *splitMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

func Split(separator string, matcher Matcher) Matcher {
	return &splitMatcher{separator: separator, matcher: matcher}
}
