package matcher

import "fmt"

type emptyMatcher struct {
}

func (m *emptyMatcher) Name() string {
	return "Empty"
}

func (m *emptyMatcher) Match(v any) (*Result, error) {
	result, err := HaveLen(0).Match(v)
	if err != nil {
		return &Result{}, err
	}

	return &Result{
		OK: result.OK,
		DescribeFailure: func() string {
			return fmt.Sprintf("%s %s %s", hint(m.Name()), _separator, v)
		},
	}, nil
}

func (m *emptyMatcher) OnMockServed() error {
	return nil
}

func (m *emptyMatcher) Spec() any {
	return _mEmpty
}

// Empty returns true if matcher value has zero length.
func Empty() Matcher {
	return &emptyMatcher{}
}
