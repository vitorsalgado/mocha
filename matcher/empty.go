package matcher

import "github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"

type emptyMatcher struct {
}

func (m *emptyMatcher) Name() string {
	return "Empty"
}

func (m *emptyMatcher) Match(v any) (*Result, error) {
	result, err := HaveLen(0).Match(v)
	if err != nil {
		return nil, err
	}

	return &Result{
		Pass:    result.Pass,
		Message: mfmt.Stringify(v),
	}, nil
}

// Empty returns true if matcher value has zero length.
func Empty() Matcher {
	return &emptyMatcher{}
}
