package matcher

import (
	"fmt"
	"strings"
)

type trimMatcher struct {
	matcher Matcher
}

func (m *trimMatcher) Match(v any) (Result, error) {
	txt, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("trim: type %T is not supported. accepted types: string", v)
	}

	result, err := m.matcher.Match(strings.TrimSpace(txt))
	if err != nil {
		return Result{}, fmt.Errorf("trim: %w", err)
	}

	if result.Pass {
		return Result{Pass: true}, nil
	}

	return Result{Message: strings.Join([]string{"Trim(", txt, ") ", result.Message}, "")}, nil
}

func (m *trimMatcher) Describe() any {
	return []any{"trim", describe(m.matcher)}
}

// Trim trims' spaces of the incoming request value before submitting it to the provided matcher.
func Trim(matcher Matcher) Matcher {
	return &trimMatcher{matcher: matcher}
}
