package matcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type allOfMatcher struct {
	matchers []Matcher
}

func (m *allOfMatcher) Match(v any) (r Result, err error) {
	var idx int
	defer func() {
		if recovery := recover(); recovery != nil {
			err = fmt.Errorf("all: matcher[%d]: panic: %v", idx, recovery)
		}
	}()

	mismatches := make([]string, 0, len(m.matchers))

	for i, matcher := range m.matchers {
		idx = i
		res, err := matcher.Match(v)
		if err != nil {
			return Result{}, fmt.Errorf("all: %d: %w", i, err)
		}

		if !res.Pass {
			mismatches = append(mismatches, res.Message)
		}
	}

	if len(mismatches) == 0 {
		return success(), nil
	}

	return Result{
		Message: strings.Join([]string{"All(", strconv.Itoa(len(m.matchers)), ")\n", mfmt.Indent(strings.Join(mismatches, "\n"))}, "")}, nil
}

// All matches when all the given matchers pass.
// Example:
//
//	All(Equal("test"), EqualIgnoreCase("test"), Contain("tes"))
func All(matchers ...Matcher) Matcher {
	return &allOfMatcher{matchers: matchers}
}
