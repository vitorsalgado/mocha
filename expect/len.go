package expect

import (
	"reflect"
)

type LenMatcher struct {
	Length int
}

func (m *LenMatcher) Name() string {
	return "Len"
}

func (m *LenMatcher) Match(v any) (Result, error) {
	value := reflect.ValueOf(v)

	return Result{
		OK: value.Len() == m.Length,
		DescribeFailure: func() string {
			return hint(m.Name(), printExpected(m.Length))
		},
	}, nil
}

func (m *LenMatcher) OnMockServed() error {
	return nil
}

// ToHaveLen returns true when matcher argument length is equal to the expected value.
func ToHaveLen(length int) Matcher {
	return &LenMatcher{Length: length}
}
