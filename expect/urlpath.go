package expect

import (
	"fmt"
	"net/url"
	"strings"
)

type URLPathMatcher struct {
	Expected string

	u string
}

func (m *URLPathMatcher) Name() string {
	return "URLPath"
}

func (m *URLPathMatcher) Match(v any) (Result, error) {
	message := func() string {
		return fmt.Sprintf(
			"%s %s %s",
			hint(m.Name(), printExpected(m.Expected)),
			_separator,
			printReceived(m.u),
		)
	}

	switch e := v.(type) {
	case *url.URL:
		m.u = e.Path
		return Result{
			OK:              strings.EqualFold(m.Expected, e.Path),
			DescribeFailure: message,
		}, nil
	case url.URL:
		m.u = e.Path
		return Result{
			OK:              strings.EqualFold(m.Expected, e.Path),
			DescribeFailure: message,
		}, nil
	case string:
		u, err := url.Parse(e)
		if err != nil {
			return Result{}, err
		}

		m.u = u.Path

		return Result{OK: strings.EqualFold(m.Expected, u.Path), DescribeFailure: message}, nil

	default:
		panic("URLPath matcher only accepts the types: *url.URL | url.URL | string")
	}
}

func (m *URLPathMatcher) OnMockServed() error {
	return nil
}

// URLPath returns true if request URL path is equal to the expected path, ignoring case.
func URLPath(expected string) Matcher {
	return &URLPathMatcher{Expected: expected}
}
