package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
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
			failed = append(failed, err.Error())

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

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "\n"))
	}

	if !ok || err != nil {
		return &Result{
			Message: mfmt.Indent(strings.Join(failed, "\n")),
			Ext:     []string{fmt.Sprintf("+%d", len(m.matchers))},
		}, err
	}

	return &Result{Pass: true}, nil
}

func (m *anyOfMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matchers...)
}

// AnyOf matches when any of the given matchers returns true.
// Example:
//
//	AnyOf(EqualTo("test"),EqualIgnoreCase("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	if len(matchers) == 0 {
		panic("[AnyOf] requires at least 1 matcher")
	}

	return &anyOfMatcher{matchers: matchers}
}
