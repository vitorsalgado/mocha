package matcher

import (
	"errors"

	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type jsonPathMatcher struct {
	path    string
	matcher Matcher
	name    string
}

func (m *jsonPathMatcher) Name() string {
	return m.name
}

func (m *jsonPathMatcher) Match(v any) (*Result, error) {
	var value any
	var err error

	if v == nil {
		value = v
	} else {
		value, err = jsonx.Reach(m.path, v)
	}

	if err != nil && !errors.Is(err, jsonx.ErrKeyNotFound) {
		return nil, err
	}

	r, err := m.matcher.Match(value)
	if err != nil {
		return &Result{}, err
	}

	if r.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{m.path, prettierName(m.matcher, r)},
		Message: r.Message,
	}, nil
}

func (m *jsonPathMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matcher)
}

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath(path string, matcher Matcher) Matcher {
	return &jsonPathMatcher{path: path, matcher: matcher, name: "JSONPath"}
}

// Field is an alias for JSONPath.
// It applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	Field("address.city", EqualTo("Santiago"))
func Field(path string, matcher Matcher) Matcher {
	return &jsonPathMatcher{path: path, matcher: matcher, name: "Field"}
}
