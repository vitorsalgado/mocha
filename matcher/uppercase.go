package matcher

import (
	"strings"
)

type upperCaseMatcher struct {
	matcher Matcher
}

func (m *upperCaseMatcher) Name() string {
	return m.matcher.Name()
}

func (m *upperCaseMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.matcher.Match(strings.ToUpper(txt))
	if err != nil {
		return nil, err
	}

	if result.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{m.matcher.Name()},
		Message: result.Message,
	}, nil
}

func (m *upperCaseMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// ToUpper upper case matcher string argument before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &upperCaseMatcher{matcher: matcher}
}
