package matcher

import (
	"fmt"
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

func (m *lowerCaseMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

func (m *lowerCaseMatcher) Spec() any {
	return []any{_mLowerCase, m.matcher.Spec()}
}

// ToLower lower case matcher string argument before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &lowerCaseMatcher{matcher: matcher}
}
