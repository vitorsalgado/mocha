package matcher

import (
	"fmt"
	"strings"
)

type allOfMatcher struct {
	matchers []Matcher
}

func (m *allOfMatcher) Name() string {
	return "AllOf"
}

func (m *allOfMatcher) Match(v any) (*Result, error) {
	ok := true
	errs := make([]string, 0)
	failed := make([]string, 0)

	for _, matcher := range m.matchers {
		result, err := matcher.Match(v)
		if err != nil {
			ok = false
			errs = append(errs, err.Error())
			failed = append(failed, result.Message())

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed, result.Message())
		}
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
			Pass:    false,
			Message: describeFailure,
		}, fmt.Errorf(strings.Join(errs, "\n"))
	}

	return &Result{Pass: ok, Message: describeFailure}, nil
}

func (m *allOfMatcher) OnMockServed() error {
	return multiOnMockServed(m.matchers...)
}

func (m *allOfMatcher) Spec() any {
	args := make([]any, len(m.matchers))

	for i, matcher := range m.matchers {
		args[i] = matcher.Spec()
	}

	return []any{_mAllOf, args}
}

// AllOf matches when all the given matchers returns true.
// Example:
//
//	AllOf(EqualTo("test"),EqualIgnoreCase("test"),ToContains("tes"))
func AllOf(matchers ...Matcher) Matcher {
	return &allOfMatcher{matchers: matchers}
}
