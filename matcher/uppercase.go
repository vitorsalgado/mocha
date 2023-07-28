package matcher

import (
	"fmt"
	"strings"
)

type upperCaseMatcher struct {
	matcher Matcher
}

func (m *upperCaseMatcher) Match(v any) (Result, error) {
	txt, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("upper: type %T is not supported. accepted types: string", v)
	}

	result, err := m.matcher.Match(strings.ToUpper(txt))
	if err != nil {
		return Result{}, fmt.Errorf("upper: %w", err)
	}

	if result.Pass {
		return Result{Pass: true}, nil
	}

	return Result{Message: strings.Join([]string{"Upper(", txt, ") ", result.Message}, "")}, nil
}

func (m *upperCaseMatcher) Describe() any {
	return []any{"uppercase", describe(m.matcher)}
}

// ToUpper uppers the case of the incoming request value before submitting it to provided matcher.
func ToUpper(matcher Matcher) Matcher {
	return &upperCaseMatcher{matcher: matcher}
}
