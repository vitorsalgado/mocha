package matcher

import (
	"errors"
	"fmt"
)

type truthyMatcher struct {
}

func (m *truthyMatcher) Name() string {
	return "Truthy"
}

func (m *truthyMatcher) Match(v any) (*Result, error) {
	b, ok := v.(bool)
	if !ok {
		return nil, errors.New("truthy matcher only works with bool values")
	}

	if !b {
		return &Result{Message: fmt.Sprintf("%s %v", hint(m.Name()), v)}, nil
	}

	return &Result{Pass: true}, nil
}

func (m *truthyMatcher) After() error {
	return nil
}

func Truthy() Matcher {
	return &truthyMatcher{}
}
