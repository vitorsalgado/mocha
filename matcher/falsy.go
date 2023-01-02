package matcher

import (
	"errors"
	"fmt"
)

type falsyMatcher struct {
}

func (m *falsyMatcher) Name() string {
	return "Falsy"
}

func (m *falsyMatcher) Match(v any) (*Result, error) {
	b, ok := v.(bool)
	if !ok {
		return nil, errors.New("falsy matcher only works with bool values")
	}

	if b {
		return &Result{Message: fmt.Sprintf("%s %v", hint(m.Name()), v)}, nil
	}

	return &Result{Pass: true}, nil
}

func (m *falsyMatcher) After() error {
	return nil
}

func Falsy() Matcher {
	return &falsyMatcher{}
}
