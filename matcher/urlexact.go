package matcher

import (
	"fmt"
	"net/url"
	"strings"
)

// TODO: complete this

type urlExactMatcher struct {
	expected string
}

func (m *urlExactMatcher) Name() string {
	return "URLPath"
}

func (m *urlExactMatcher) Match(v any) (*Result, error) {
	message := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.expected)),
			_separator,
			printReceived(m),
		)
	}

	switch e := v.(type) {
	case *url.URL:
		return &Result{
			OK:              strings.EqualFold(m.expected, e.String()),
			DescribeFailure: message,
		}, nil
	case string:
		u, err := url.Parse(e)
		if err != nil {
			return &Result{}, err
		}

		return &Result{OK: strings.EqualFold(m.expected, u.String()), DescribeFailure: message}, nil

	default:
		panic("URLPath matcher only accepts the types: *url.URL | url.URL | string")
	}
}

func (m *urlExactMatcher) OnMockServed() error {
	return nil
}

// URLExact matches the entire URL converted to string.
func URLExact(expected string) Matcher {
	return &urlPathMatcher{expected: expected}
}
