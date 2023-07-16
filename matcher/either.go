package matcher

import (
	"fmt"
	"strings"
)

type eitherMatcher struct {
	first  Matcher
	second Matcher
}

func (m *eitherMatcher) Match(v any) (Result, error) {
	r1, err := m.first.Match(v)
	if err != nil {
		return Result{}, fmt.Errorf("either: first: %w", err)
	}

	r2, err := m.second.Match(v)
	if err != nil {
		return Result{}, fmt.Errorf("either: second: %w", err)
	}

	if r1.Pass || r2.Pass {
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

	return Result{Message: strings.Join([]string{"Either() ", message}, "")}, nil
}

// Either passes when any of the two given matchers pass.
func Either(first Matcher, second Matcher) Matcher {
	return &eitherMatcher{first: first, second: second}
}
