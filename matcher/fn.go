package matcher

import "fmt"

type funcMatcher struct {
	Func func(v any) (bool, error)
}

func (m *funcMatcher) Name() string {
	return "Func"
}

func (m *funcMatcher) Match(v any) (*Result, error) {
	r, err := m.Func(v)
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		OK: r,
		DescribeFailure: func() string {
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
	return &funcMatcher{Func: fn}
}
