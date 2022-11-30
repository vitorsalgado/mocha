package matcher

import "fmt"

type EmptyMatcher struct {
}

func (m *EmptyMatcher) Name() string {
	return "Empty"
}

func (m *EmptyMatcher) Match(v any) (Result, error) {
	result, err := HaveLen(0).Match(v)
	if err != nil {
		return Result{}, err
	}

	return Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %s", hint(m.Name()), _separator, v)
		},
	}, nil
}

func (m *EmptyMatcher) OnMockServed() error {
	return nil
}

// Empty returns true if matcher value has zero length.
func Empty() Matcher {
	return &EmptyMatcher{}
}
