package expect

import "fmt"

type FuncMatcher struct {
	Func func(v any) (bool, error)
}

func (m *FuncMatcher) Name() string {
	return "Func"
}

func (m *FuncMatcher) Match(v any) (bool, error) {
	return m.Func(v)
}

func (m *FuncMatcher) DescribeFailure(v any) string {
	return fmt.Sprintf("custom matcher function did not match. value: %v", v)
}

func (m *FuncMatcher) OnMockServed() error {
	return nil
}

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any) (bool, error)) Matcher {
	return &FuncMatcher{Func: fn}
}
