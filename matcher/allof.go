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
			failed = append(failed, err.Error())

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed, result.Message)
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "\n"))
	}

	if !ok || err != nil {
		return &Result{
			Message: indent(strings.Join(failed, "\n")),
			Ext:     []string{fmt.Sprintf("+%d", len(m.matchers))},
		}, err
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
