package matcher

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type equalStrictMatcher struct {
	expected any
}

func (m *equalStrictMatcher) Match(v any) (Result, error) {
	if reflect.DeepEqual(m.expected, v) {
		return Result{Pass: true}, nil
	}

	return mismatch(strings.Join([]string{"StrictEq(", mfmt.Stringify(m.expected), ") Got: ", mfmt.Stringify(v)}, "")), nil
}

func (m *equalStrictMatcher) Describe() any {
	return []any{"eqs", m.expected}
}

// StrictEqual strictly compares the expected value with incoming request values, considering value and type.
func StrictEqual(expected any) Matcher {
	return &equalStrictMatcher{expected: expected}
}

// StrictEqualf strictly compares the expected value with incoming request values, considering value and type.
// This is short-hand to format the expected value.
func StrictEqualf(format string, a ...any) Matcher {
	return &equalStrictMatcher{expected: fmt.Sprintf(format, a...)}
}

// Eqs strictly compares the expected value with incoming request values, considering value and type.
func Eqs(expected any) Matcher {
	return &equalStrictMatcher{expected: expected}
}
