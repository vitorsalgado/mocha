package matcher

import (
	"reflect"
)

type lenMatcher struct {
	Length int
}

func (m *lenMatcher) Name() string {
	return "Len"
}

func (m *lenMatcher) Match(v any) (*Result, error) {
	value := reflect.ValueOf(v)

	return &Result{
		OK: value.Len() == m.Length,
		DescribeFailure: func() string {
			return hint(m.Name(), printExpected(m.Length))
		},
	}, nil
}

func (m *lenMatcher) OnMockServed() error {
	return nil
}

// HaveLen returns true when matcher argument length is equal to the expected value.
func HaveLen(length int) Matcher {
	return &lenMatcher{Length: length}
}
