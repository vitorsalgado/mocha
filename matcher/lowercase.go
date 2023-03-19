package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type lowerCaseMatcher struct {
	matcher Matcher
}

func (m *lowerCaseMatcher) Name() string {
	return "ToLower"
}

func (m *lowerCaseMatcher) Match(v any) (*Result, error) {
	txt, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("type %v is not supported. accepted types: string", reflect.TypeOf(v))
	}

	result, err := m.matcher.Match(strings.ToLower(txt))
	if err != nil {
		return nil, err
	}

	if result.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{txt, prettierName(m.matcher, result)},
		Message: result.Message,
	}, nil
}

func (m *lowerCaseMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// ToLower lowers the incoming request value case value before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &lowerCaseMatcher{matcher: matcher}
}
