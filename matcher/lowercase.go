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
	txt := v.(string)
	result, err := m.matcher.Match(strings.ToLower(txt))
	if err != nil {
		return &Result{}, err
	}

	if result.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{stringify(txt)},
		Message: result.Message,
	}, nil
}

func (m *lowerCaseMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// ToLower lower case matcher string argument before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &lowerCaseMatcher{matcher: matcher}
}
