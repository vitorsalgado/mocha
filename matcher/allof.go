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
			failed = append(failed, result.Message)

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed, result.Message)
		}
	}

	if len(errs) > 0 {
		return &Result{
			Pass: false,
			Message: fmt.Sprintf(
				"%s\n%s",
				hint(m.Name(), fmt.Sprintf("+%d", len(m.matchers))),
				indent(strings.Join(failed, "\n")),
			),
		}, fmt.Errorf(strings.Join(errs, "\n"))
	}

	if !ok {
		return &Result{Message: fmt.Sprintf(
			"%s\n%s",
			hint(m.Name(), fmt.Sprintf("+%d", len(m.matchers))),
			indent(strings.Join(failed, "\n")),
		)}, nil
	}

	return &Result{Pass: true}, nil
}

func (m *allOfMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matchers...)
}

// AllOf matches when all the given matchers returns true.
// Example:
//
//	AllOf(EqualTo("test"),EqualIgnoreCase("test"),ToContains("tes"))
func AllOf(matchers ...Matcher) Matcher {
	if len(matchers) == 0 {
		panic("[AllOf] requires at least 1 matcher")
	}

	return &allOfMatcher{matchers: matchers}
}
