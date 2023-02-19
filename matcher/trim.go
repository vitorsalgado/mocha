package matcher

import (
	"strings"
)

type trimMatcher struct {
	matcher Matcher
}

func (m *trimMatcher) Name() string {
	return "Trim"
}

func (m *trimMatcher) Match(v any) (*Result, error) {
	txt := v.(string)
	result, err := m.matcher.Match(strings.TrimSpace(txt))
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

func (m *trimMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// Trim trims' spaces of matcher argument before submitting it to the given matcher.
func Trim(matcher Matcher) Matcher {
	return &trimMatcher{matcher: matcher}
}
