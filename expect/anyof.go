package expect

import (
	"fmt"
	"strings"
)

type AnyOfMatcher struct {
	Matchers []Matcher
}

func (m *AnyOfMatcher) Name() string {
	return "AnyOf"
}

func (m *AnyOfMatcher) Match(v any) (Result, error) {
	ok := false
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

			continue
		}

		ok = true
		break
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

func (m *AnyOfMatcher) OnMockServed() error {
	return multiOnMockServed(m.Matchers...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	return &AnyOfMatcher{Matchers: matchers}
}
