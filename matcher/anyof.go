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
			failed = append(failed, result.Message)

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed, result.Message)

			continue
		}

		ok = true
		break
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

	return &Result{Pass: ok}, nil
}

func (m *anyOfMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matchers...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),EqualIgnoreCase("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	return &anyOfMatcher{matchers: matchers}
}
