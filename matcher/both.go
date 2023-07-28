package matcher

import (
	"fmt"
	"strings"
)

type bothMatcher struct {
	first  Matcher
	second Matcher
}

func (m *bothMatcher) Name() string {
	return "Both"
}

func (m *bothMatcher) Match(value any) (Result, error) {
	r1, err := m.first.Match(value)
	if err != nil {
		return Result{}, fmt.Errorf("both: first: %w", err)
	}

	r2, err := m.second.Match(value)
	if err != nil {
		return Result{}, fmt.Errorf("both: first: %w", err)
	}

	if r1.Pass && r2.Pass {
		return Result{Pass: true}, nil
	}

	message := ""
	if !r1.Pass {
		message = r1.Message
	}

	if !r2.Pass {
		if !r1.Pass {
			message += " - "
		}
		message += r2.Message
	}

	return Result{Message: strings.Join([]string{"Both() ", message}, "")}, nil
}

func (m *bothMatcher) Describe() any {
	return []any{"both", describe(m.first), describe(m.second)}
}

// Both passes when both the given matchers pass.
func Both(first Matcher, second Matcher) Matcher {
	return &bothMatcher{first: first, second: second}
}
