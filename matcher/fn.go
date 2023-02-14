package matcher

import "fmt"

type funcMatcher struct {
	fn func(v any) (bool, error)
}

func (m *funcMatcher) Name() string {
	return "fn"
}

func (m *funcMatcher) Match(v any) (*Result, error) {
	pass, err := m.fn(v)
	if err != nil {
		return nil, err
	}

	if pass {
		return &Result{Pass: true}, nil
	}

	return &Result{
		Ext:     []string{stringify(v)},
		Message: fmt.Sprintf("Received: %v", v),
	}, nil
}

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any) (bool, error)) Matcher {
	return &funcMatcher{fn: fn}
}
