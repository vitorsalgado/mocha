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
	if value.Len() == m.length {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: stringify(m.length)}, nil
}

// HaveLen returns true when matcher argument length is equal to the items value.
func HaveLen(length int) Matcher {
	return &lenMatcher{length: length}
}
