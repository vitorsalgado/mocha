package expect

import (
	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type JSONPathMatcher struct {
	Path    string
	Matcher Matcher
}

func (m *JSONPathMatcher) Name() string {
	return "JSONPath"
}

func (m *JSONPathMatcher) Match(v any) (Result, error) {
	var value any
	var err error

	if v == nil {
		value = v
	} else {
		value, err = jsonx.Reach(m.Path, v)
	}

	if err != nil {
		return mismatch(nil), err
	}

	r, err := m.Matcher.Match(value)
	if err != nil {
		return Result{}, err
	}

	return Result{OK: r.OK, DescribeFailure: func() string {
		return hint(m.Name(), printExpected(m.Path), r.DescribeFailure())
	}}, nil
}

func (m *JSONPathMatcher) OnMockServed() error {
	return m.Matcher.OnMockServed()
}

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath(path string, matcher Matcher) Matcher {
	return &JSONPathMatcher{Path: path, Matcher: matcher}
}
