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
	return "Any"
}

func (m *anyOfMatcher) Match(v any) (r *Result, e error) {
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

	ok := false
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
			Message: "\n" + mfmt.Indent(strings.Join(failed, "\n")),
			Ext:     []string{fmt.Sprintf("+%d", len(m.matchers))},
		}, err
	}

	return &Result{Pass: true}, nil
}

func (m *anyOfMatcher) AfterMockServed() error {
	return runAfterMockServed(m.matchers...)
}

// Any matches when any of the given matchers pass.
// Example:
//
//	Any(Equal("test"), EqualIgnoreCase("TEST"), Contain("tes"))
func Any(matchers ...Matcher) Matcher {
	if len(matchers) == 0 {
		panic("matcher: [Any] requires at least 1 matcher")
	}

	return &anyOfMatcher{matchers: matchers}
}
