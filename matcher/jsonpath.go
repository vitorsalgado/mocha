package matcher

import (
	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type jsonPathMatcher struct {
	path    string
	matcher Matcher
}

func (m *jsonPathMatcher) Name() string {
	return "JSONPath"
}

func (m *jsonPathMatcher) Match(v any) (*Result, error) {
	var value any
	var err error

	if v == nil {
		value = v
	} else {
		value, err = jsonx.Reach(m.path, v)
	}

	if err != nil {
		return nil, err
	}

	r, err := m.matcher.Match(value)
	if err != nil {
		return &Result{}, err
	}

	return &Result{Pass: r.Pass, Message: func() string {
		return hint(m.Name(), printExpected(m.path), r.Message())
	}}, nil
}

func (m *jsonPathMatcher) OnMockServed() error {
	return m.matcher.OnMockServed()
}

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath(path string, matcher Matcher) Matcher {
	return &jsonPathMatcher{path: path, matcher: matcher}
}
