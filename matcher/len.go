package matcher

import (
	"reflect"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type lenMatcher struct {
	length int
}

func (m *lenMatcher) Name() string {
	return "Len"
}

func (m *lenMatcher) Match(v any) (*Result, error) {
	value := reflect.ValueOf(v)
	if value.Len() == m.length {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: mfmt.Stringify(m.length)}, nil
}

// Len passes when the expected value length is equal to the incoming request value.
func Len(length int) Matcher {
	return &lenMatcher{length: length}
}
