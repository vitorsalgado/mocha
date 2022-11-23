package expect

import (
	"fmt"

	"github.com/vitorsalgado/mocha/v3/internal/jsonx"
)

type JSONPathMatcher struct {
	Path    string
	Matcher Matcher
}

func (m *JSONPathMatcher) Name() string {
	return "JSONPath"
}

func (m *JSONPathMatcher) Match(v any) (bool, error) {
	value, err := jsonx.Reach(m.Path, v)
	if err != nil || value == nil {
		return false, err
	}

	return m.Matcher.Match(value)
}

func (m *JSONPathMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("matcher %s applied on json field %s did not match", m.Matcher.Name(), m.Path)
}

func (m *JSONPathMatcher) OnMockServed() error {
	return nil
}

// JSONPath applies the provided matcher to the JSON field value in the given path.
// Example:
//
//	JSONPath("address.city", EqualTo("Santiago"))
func JSONPath(path string, matcher Matcher) Matcher {
	return &JSONPathMatcher{Path: path, Matcher: matcher}
}
