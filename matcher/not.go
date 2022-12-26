package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/types"
)

type notMatcher struct {
	matcher Matcher
}

func (m *notMatcher) Name() string {
	return "Not"
}

func (m *notMatcher) Match(v any) (*Result, error) {
	result, err := m.matcher.Match(v)
	if err != nil {
		return nil, err
	}

	return &Result{
		Pass: !result.Pass,
		Message: func() string {
			return fmt.Sprintf(
				"%s ! %s",
				hint(m.Name(), m.matcher.Name()),
				result.Message(),
			)
		},
	}, nil
}

func (m *notMatcher) AfterMockSent() error {
	return m.matcher.AfterMockSent()
}

func (m *notMatcher) Raw() types.RawValue {
	return types.RawValue{_mNot, m.matcher.Raw()}
}

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	return &notMatcher{matcher: matcher}
}
