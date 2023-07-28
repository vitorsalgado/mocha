package matcher

import (
	"fmt"
	"strings"
)

type splitMatcher struct {
	separator string
	matcher   Matcher
}

func (m *splitMatcher) Match(v any) (Result, error) {
	txt, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("split: type %T is not supported. accepted types: string", v)
	}

	result, err := m.matcher.Match(strings.Split(txt, m.separator))
	if err != nil {
		return Result{}, fmt.Errorf("split: %w", err)
	}

	if result.Pass {
		return Result{Pass: true}, err
	}

	return Result{Message: strings.Join([]string{"Split(", m.separator, ") ", result.Message}, "")}, nil
}

func (m *splitMatcher) Describe() any {
	return []any{"split", m.separator, describe(m.matcher)}
}

// Split splits the incoming request value text using the separator parameter and
// test each element with the provided matcher.
func Split(separator string, matcher Matcher) Matcher {
	return &splitMatcher{separator: separator, matcher: matcher}
}
