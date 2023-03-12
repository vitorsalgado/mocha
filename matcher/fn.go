package matcher

import (
	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type funcMatcher struct {
	fn func(v any) (bool, error)
}

func (m *funcMatcher) Name() string {
	return "Func"
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
		Ext:     []string{mfmt.Stringify(v)},
		Message: "predicate evaluated to false",
	}, nil
}

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any) (bool, error)) Matcher {
	return &funcMatcher{fn: fn}
}
