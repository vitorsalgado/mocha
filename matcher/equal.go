package matcher

import (
	"fmt"
	"reflect"

	"github.com/vitorsalgado/mocha/v3/types"
)

type equalMatcher struct {
	expected any
}

func (m *equalMatcher) Name() string {
	return "Equal"
}

func (m *equalMatcher) Match(v any) (*Result, error) {
	return &Result{
		Pass: reflect.DeepEqual(m.expected, v),
		Message: func() string {
			return fmt.Sprintf("%s %s %v",
				hint(m.Name(), printExpected(m.expected)),
				_separator,
				printReceived(v),
			)
		},
	}, nil
}

func (m *equalMatcher) AfterMockSent() error {
	return nil
}

func (m *equalMatcher) Raw() types.RawValue {
	return types.RawValue{_mEqualTo, m.expected}
}

// Equal returns true if matcher value is equal to the given parameter value.
func Equal(expected any) Matcher {
	return &equalMatcher{expected: expected}
}

// Equalf returns true if matcher value is equal to the given parameter value.
func Equalf(format string, a ...any) Matcher {
	return &equalMatcher{expected: fmt.Sprintf(format, a...)}
}
