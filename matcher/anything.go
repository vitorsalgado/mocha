package matcher

import "github.com/vitorsalgado/mocha/v3/types"

type anythingMatcher struct {
}

func (m *anythingMatcher) Name() string {
	return "Anything"
}

func (m *anythingMatcher) Match(v any) (*Result, error) {
	return &Result{Pass: true}, nil
}

func (m *anythingMatcher) AfterMockSent() error {
	return nil
}

func (m *anythingMatcher) Raw() types.RawValue {
	return nil
}

func Anything() Matcher {
	return &anythingMatcher{}
}
