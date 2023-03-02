package matcher

import "github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"

type emptyMatcher struct {
}

func (m *emptyMatcher) Name() string {
	return "Empty"
}

func (m *emptyMatcher) Match(v any) (*Result, error) {
	if v == nil {
		return &Result{Pass: true}, nil
	}

	result, err := HasLen(0).Match(v)
	if err != nil {
		return nil, err
	}

	return &Result{
		Pass:    result.Pass,
		Message: mfmt.Stringify(v),
	}, nil
}

// Empty passes if the incoming request value is empty or has zero value.
func Empty() Matcher {
	return &emptyMatcher{}
}
