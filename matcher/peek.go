package matcher

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/types"
)

type peekMatcher struct {
	matcher Matcher
	action  func(v any) error
}

func (m *peekMatcher) Name() string {
	return fmt.Sprintf("Peek(%s)", m.matcher.Name())
}

func (m *peekMatcher) Match(v any) (*Result, error) {
	err := m.action(v)
	if err != nil {
		return nil, err
	}

	return m.matcher.Match(v)
}

func (m *peekMatcher) AfterMockSent() error {
	return m.matcher.AfterMockSent()
}

func (m *peekMatcher) Raw() types.RawValue {
	return nil
}

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	return &peekMatcher{matcher: matcher, action: action}
}
