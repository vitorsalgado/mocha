package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type emptyMatcher struct {
}

func (m *emptyMatcher) Match(v any) (Result, error) {
	if v == nil {
		return Result{Pass: true}, nil
	}

	result, err := Len(0).Match(v)
	if err != nil {
		return Result{}, fmt.Errorf("empty: %w", err)
	}

	return Result{result.Pass, strings.Join([]string{"Empty(", mfmt.Stringify(v), ")"}, "")}, nil
}

func (m *emptyMatcher) Describe() any {
	return []any{"empty"}
}

// Empty passes if the incoming request value is empty or has zero value.
func Empty() Matcher {
	return &emptyMatcher{}
}
