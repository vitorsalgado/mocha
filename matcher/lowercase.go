package matcher

import (
	"strings"
)

type lowerCaseMatcher struct {
	matcher Matcher
}

func (m *lowerCaseMatcher) Name() string {
	return "ToLower"
}

func (m *lowerCaseMatcher) Match(v any) (*Result, error) {
	// TODO: check for cast errors
	txt := v.(string)
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
