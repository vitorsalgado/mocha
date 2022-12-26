package matcher

import "fmt"

type funcMatcher struct {
	fn func(v any) (bool, error)
}

func (m *funcMatcher) Name() string {
	return "fn"
}

func (m *funcMatcher) Match(v any) (*Result, error) {
	r, err := m.fn(v)
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		Pass: r,
		Message: func() string {
			return fmt.Sprintf(
				"%s %s %v",
				hint(m.Name()),
				_separator,
				v,
			)
		},
	}, nil
}

func (m *funcMatcher) OnMockServed() error {
	return nil
}

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any) (bool, error)) Matcher {
	return &funcMatcher{fn: fn}
}
