package expect

import (
	"fmt"
	"strings"
)

type AllOfMatcher struct {
	Matchers []Matcher
}

func (m *AllOfMatcher) Name() string {
	return "AllOf"
}

func (m *AllOfMatcher) Match(v any) (Result, error) {
	ok := true
	errs := make([]string, 0)
	failed := make([]string, 0)

	for _, matcher := range m.Matchers {
		result, err := matcher.Match(v)
		if err != nil {
			ok = false
			errs = append(errs, err.Error())
			failed = append(failed, result.DescribeFailure())

			continue
		}

		if !result.OK {
			ok = false
			failed = append(failed, result.DescribeFailure())
		}
	}

	describeFailure := func() string {
		return fmt.Sprintf(
			"%s\n%s",
			hint(m.Name(), fmt.Sprintf("+%d", len(m.Matchers))),
			indent(strings.Join(failed, "\n")),
		)
	}

	if len(errs) > 0 {
		return Result{
			OK:              false,
			DescribeFailure: describeFailure,
		}, fmt.Errorf(strings.Join(errs, "\n"))
	}

	return Result{OK: ok, DescribeFailure: describeFailure}, nil
}

func (m *AllOfMatcher) OnMockServed() error {
	return multiOnMockServed(m.Matchers...)
}

// AllOf matches when all the given matchers returns true.
// Example:
//
//	AllOf(EqualTo("test"),ToEqualFold("test"),ToContains("tes"))
func AllOf(matchers ...Matcher) Matcher {
	return &AllOfMatcher{Matchers: matchers}
}
