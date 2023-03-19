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
	return "Split"
}

func (m *splitMatcher) Match(v any) (*Result, error) {
	txt, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("type %v is not supported. accepted types: string", reflect.TypeOf(v))
	}

	result, err := m.matcher.Match(strings.Split(txt, m.separator))
	if err != nil {
		return nil, err
	}

	if result.Pass {
		return &Result{Pass: true}, err
	}

	return &Result{
		Ext:     []string{m.separator, prettierName(m.matcher, result)},
		Message: result.Message,
	}, nil
}

func (m *splitMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// Split splits the incoming request value text using the separator parameter and
// test each element with the provided matcher.
func Split(separator string, matcher Matcher) Matcher {
	return &splitMatcher{separator: separator, matcher: matcher}
}
