package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type equalMatcher struct {
	expected any
}

func (m *equalMatcher) Match(v any) (Result, error) {
	if equalValues(m.expected, v, true) {
		return success(), nil
	}

	return mismatch(strings.Join([]string{"Equal(", mfmt.Stringify(m.expected), ") Got: ", mfmt.Stringify(v)}, "")), nil
}

func (m *equalMatcher) Describe() any {
	return []any{"eq", m.expected}
}

// Equal asserts that the given expectation is equal to the incoming request value.
// It considers equivalent value. Eg.: float64(10) is equal to int(10).
func Equal(expected any) Matcher {
	return &equalMatcher{expected: expected}
}

// Equalf returns true if the matcher value is equal to the given parameter value.
// It considers equivalent value. Eg.: float64(10) is equal to int(10).
// This is short-hand to format the expected value.
func Equalf(format string, a ...any) Matcher {
	return &equalMatcher{expected: fmt.Sprintf(format, a...)}
}

// Eq is an alias to Equal.
// It asserts that the given expectation is equal to the incoming request value.
// It considers equivalent value. Eg.: float64(10) is equal to int(10).
func Eq(expected any) Matcher {
	return &equalMatcher{expected: expected}
}
