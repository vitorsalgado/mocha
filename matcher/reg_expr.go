package matcher

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type regExpMatcher struct {
	expression string
	rg         *regexp.Regexp
	mu         sync.Mutex
}

func (m *regExpMatcher) Match(v any) (Result, error) {
	txt, ok := v.(string)
	if !ok {
		return Result{}, fmt.Errorf("regexp: it only works with string values. got: %T", v)
	}

	if m.rg == nil {
		m.mu.Lock()
		if m.rg == nil {
			r, err := regexp.Compile(m.expression)
			if err != nil {
				return Result{}, fmt.Errorf("regexp: error compiling expression %s. %w", m.expression, err)
			}

			m.rg = r
		}
		m.mu.Unlock()
	}

	match := m.rg.Match([]byte(txt))
	if match {
		return success(), nil
	}

	return Result{Message: strings.Join([]string{"Match(", m.expression, ") Expression did not match. Got: ", txt}, "")}, nil
}

// Matches passes when the given regular expression matches the incoming request value.
// It accepts a string, regexp.Regexp or *regexp.Regexp.
func Matches(expression string) Matcher {
	return &regExpMatcher{expression: expression}
}
