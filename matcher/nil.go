package matcher

import "github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"

type nilMatcher struct {
}

func (m *nilMatcher) Name() string {
	return "Nil"
}

func (m *nilMatcher) Match(v any) (*Result, error) {
	if v == nil {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: mfmt.PrintReceived(v)}, nil
}

func Nil() Matcher {
	return &nilMatcher{}
}
