package matcher

import (
	"reflect"
)

type lenMatcher struct {
	length int
}

func (m *lenMatcher) Name() string {
	return "Len"
}

func (m *lenMatcher) Match(v any) (*Result, error) {
	value := reflect.ValueOf(v)

	return &Result{
		Pass: value.Len() == m.length,
		Message: func() string {
			return hint(m.Name(), printExpected(m.length))
		},
	}, nil
}

func (m *lenMatcher) OnMockServed() error {
	return nil
}

func (m *lenMatcher) Spec() any {
	return []any{_mLen, m.length}
}

// HaveLen returns true when matcher argument length is equal to the items value.
func HaveLen(length int) Matcher {
	return &lenMatcher{length: length}
}
