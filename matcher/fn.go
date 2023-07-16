package matcher

import "fmt"

type funcMatcher struct {
	fn func(v any) (bool, error)
}

func (m *funcMatcher) Match(v any) (Result, error) {
	pass, err := m.fn(v)
	if err != nil {
		return Result{}, fmt.Errorf("fn: matcher: %w", err)
	}

	if pass {
		return Result{Pass: true}, nil
	}

	return Result{
		Message: "F() Wrapped predicate did not match",
	}, nil
}

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any) (bool, error)) Matcher {
	return &funcMatcher{fn: fn}
}
