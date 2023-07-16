package matcher

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vitorsalgado/mocha/v3/matcher/internal/mfmt"
)

type anyOfMatcher struct {
	matchers []Matcher
}

func (m *anyOfMatcher) Match(v any) (r Result, e error) {
	var idx int
	defer func() {
		if recovery := recover(); recovery != nil {
			e = fmt.Errorf("any: matcher[%d]: panic: %v", idx, recovery)
		}
	}()

	mismatches := make([]string, 0, len(m.matchers))

	for i, matcher := range m.matchers {
		idx = i
		result, err := matcher.Match(v)
		if err != nil {
			return Result{}, fmt.Errorf("any: %d: %w", i, err)
		}

		if !result.Pass {
			mismatches = append(mismatches, result.Message)
			continue
		}

		return success(), nil
	}

	return Result{Message: strings.Join([]string{"Any(", strconv.Itoa(len(m.matchers)), ")\n", mfmt.Indent(strings.Join(mismatches, "\n"))}, "")}, nil
}

// Any matches when any of the given matchers pass.
// Example:
//
//	Any(Equal("test"), EqualIgnoreCase("TEST"), Contain("tes"))
func Any(matchers ...Matcher) Matcher {
	return &anyOfMatcher{matchers: matchers}
}
