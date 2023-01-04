package matcher

import "fmt"

type nilMatcher struct {
}

func (m *nilMatcher) Name() string {
	return "Nil"
}

func (m *nilMatcher) Match(v any) (*Result, error) {
	if v == nil {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: fmt.Sprintf("%s %s %v",
		hint(m.Name()),
		_separator,
		printReceived(v)),
	}, nil
}

func Nil() Matcher {
	return &nilMatcher{}
}
