package expect

import (
	"fmt"
	"net/url"
	"strings"
)

type URLPathMatcher struct {
	Expected string
}

func (m *URLPathMatcher) Name() string {
	return "URLPath"
}

func (m *URLPathMatcher) Match(v any) (bool, error) {
	switch e := v.(type) {
	case *url.URL:
		return strings.EqualFold(m.Expected, e.Path), nil
	case url.URL:
		return strings.EqualFold(m.Expected, e.Path), nil
	case string:
		u, err := url.Parse(e)
		if err != nil {
			return false, err
		}

		return strings.EqualFold(m.Expected, u.Path), nil

	default:
		panic("URLPath matcher only accepts the types: *url.URL | url.URL | string")
	}
}

func (m *URLPathMatcher) DescribeFailure(_ any) string {
	return fmt.Sprintf("url does not have the expected path %s", m.Expected)
}

func (m *URLPathMatcher) OnMockServed() error {
	return nil
}

// URLPath returns true if request URL path is equal to the expected path, ignoring case.
func URLPath(expected string) Matcher {
	return &URLPathMatcher{Expected: expected}
}
