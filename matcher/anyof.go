package matcher

import (
	"fmt"
	"strings"
)

type anyOfMatcher struct {
	matchers []Matcher
}

func (m *anyOfMatcher) Name() string {
	return "AnyOf"
}

func (m *anyOfMatcher) Match(v any) (*Result, error) {
	ok := false
	errs := make([]string, 0)
	failed := make([]string, 0)

	for _, matcher := range m.matchers {
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
			hint(m.Name(), fmt.Sprintf("+%d", len(m.matchers))),
			indent(strings.Join(failed, "\n")),
		)
	}

	if len(errs) > 0 {
		return &Result{
			OK:              false,
			DescribeFailure: describeFailure,
		}, fmt.Errorf(strings.Join(errs, "\n"))
	}

	return &Result{OK: ok, DescribeFailure: describeFailure}, nil
}

func (m *anyOfMatcher) OnMockServed() error {
	return multiOnMockServed(m.matchers...)
}

func (m *anyOfMatcher) Spec() any {
	args := make([]any, len(m.matchers))

	for i, matcher := range m.matchers {
		args[i] = matcher.Spec()
	}

	return []any{_mAnyOf, args}
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),EqualIgnoreCase("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	return &anyOfMatcher{matchers: matchers}
}
