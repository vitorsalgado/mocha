package matcher

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type allOfMatcher struct {
	matchers []Matcher
}

func (m *allOfMatcher) Name() string {
	return "All"
}

func (m *allOfMatcher) Match(v any) (r *Result, e error) {
	var idx int
	var cur Matcher

	defer func() {
		if recovery := recover(); recovery != nil {
			n := ""
			if cur != nil {
				n = cur.Name()
			}

			r = nil
			e = fmt.Errorf("%v, matcher=%s, index=%d", recovery, n, idx)
		}
	}()

	ok := true
	errs := make([]string, 0)
	failed := make([]string, 0)

	for i, matcher := range m.matchers {
		idx = i
		cur = matcher

		result, err := matcher.Match(v)
		if err != nil {
			ok = false
			errs = append(errs, err.Error())
			failed = append(failed, err.Error())

			continue
		}

		if !result.Pass {
			ok = false
			failed = append(failed,
				fmt.Sprintf("%s(%s) %s", matcher.Name(), strings.Join(result.Ext, ", "), result.Message))
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf(strings.Join(errs, "\n"))
	}

	if !ok || err != nil {
		return &Result{
			Message: "\n" + mfmt.Indent(strings.Join(failed, "\n")),
			Ext:     []string{fmt.Sprintf("+%d", len(m.matchers))},
		}, err
	}

	return &Result{Pass: true}, nil
}

func (m *allOfMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matchers...)
}

// All matches when all the given matchers returns true.
// Example:
//
//	All(EqualTo("test"),EqualIgnoreCase("test"),ToContains("tes"))
func All(matchers ...Matcher) Matcher {
	if len(matchers) == 0 {
		panic("[All] requires at least 1 matcher")
	}

	return &allOfMatcher{matchers: matchers}
}
