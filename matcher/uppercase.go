package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type upperCaseMatcher struct {
	matcher Matcher
}

func (m *upperCaseMatcher) Name() string {
	return "ToUpper"
}

func (m *upperCaseMatcher) Match(v any) (*Result, error) {
	txt, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("type %v is not supported. accepted types: string", reflect.TypeOf(v))
	}

	result, err := m.matcher.Match(strings.ToUpper(txt))
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

func (m *upperCaseMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// ToUpper uppers the case of the incoming request value before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &upperCaseMatcher{matcher: matcher}
}
