package matcher

import (
	"fmt"
	"strings"
)

type lowerCaseMatcher struct {
	matcher Matcher
}

func (m *lowerCaseMatcher) Match(v any) (Result, error) {
	txt, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("lower: type %T is not supported. accepted types: string", v)
	}

	result, err := m.matcher.Match(strings.ToLower(txt))
	if err != nil {
		return Result{}, fmt.Errorf("lower: %w", err)
	}

	if result.Pass {
		return Result{Pass: true}, nil
	}

	return Result{
		Message: "Lower() " + result.Message,
	}, nil
}

func (m *lowerCaseMatcher) Describe() any {
	return []any{"lowercase", describe(m.matcher)}
}

// ToLower lowers the incoming request value case value before submitting it to provided matcher.
func ToLower(matcher Matcher) Matcher {
	return &lowerCaseMatcher{matcher: matcher}
}
