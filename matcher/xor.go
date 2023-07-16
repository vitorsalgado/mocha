package matcher

import (
	"fmt"
	"strings"
)

type xorMatcher struct {
	first  Matcher
	second Matcher
}

func (m *xorMatcher) Match(v any) (Result, error) {
	a, err := m.first.Match(v)
	if err != nil {
		return Result{}, fmt.Errorf("xor: first: %w", err)
	}

	b, err := m.second.Match(v)
	if err != nil {
		return Result{}, fmt.Errorf("xor: second: %w", err)
	}

	if a.Pass != b.Pass {
		return Result{Pass: true}, nil
	}

	message := ""
	if !a.Pass {
		message = a.Message
	}
	if !b.Pass {
		if !b.Pass {
			message += " - "
		}
		message += b.Message
	}

	return Result{Message: strings.Join([]string{"XOR()", message}, " ")}, nil
}

// XOR is an exclusive OR matcher.
func XOR(first Matcher, second Matcher) Matcher {
	return &xorMatcher{first: first, second: second}
}
