package matcher

import (
	"reflect"

	"github.com/vitorsalgado/mocha/v3/types"
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

func (m *lenMatcher) AfterMockSent() error {
	return nil
}

func (m *lenMatcher) Raw() types.RawValue {
	return types.RawValue{_mLen, m.length}
}

// HaveLen returns true when matcher argument length is equal to the items value.
func HaveLen(length int) Matcher {
	return &lenMatcher{length: length}
}
