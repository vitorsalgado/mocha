package matcher

import (
	"fmt"
	"net/url"
)

type urlPathMatcher struct {
	matcher Matcher
}

func (m *urlPathMatcher) Name() string {
	return "URLPath"
}

func (m *urlPathMatcher) Match(v any) (*Result, error) {
	var value any

	switch e := v.(type) {
	case *url.URL:
		value = e.Path
	case string:
		u, err := url.Parse(e)
		if err != nil {
			return nil, err
		}

		value = u.Path
	case fmt.Stringer:
		u, err := url.Parse(e.String())
		if err != nil {
			return nil, err
		}

		value = u.Path
	default:
		panic("URLPath matcher only accepts the types: *url.URL | url.URL | string")
	}

	res, err := m.matcher.Match(value)
	if err != nil {
		return nil, err
	}

	if res.Pass {
		return &Result{Pass: true}, nil
	}

	return &Result{Message: res.Message, Ext: []string{prettierName(m.matcher, res)}}, nil
}

// URLPath compares the URL path with the expected value and matches if they are equal.
// Comparison is case-insensitive.
func URLPath(expected string) Matcher {
	return URLPathMatch(EqualIgnoreCase(expected))
}

// URLPathf compares the URL path with the expected value and matches if they are equal.
// The expected value will be formatted according to the given format specifier.
// Comparison is case-insensitive.
func URLPathf(format string, a ...any) Matcher {
	return URLPathMatch(EqualIgnoreCase(fmt.Sprintf(format, a...)))
}

// URLPathMatch applies the provided matcher to the URL path.
func URLPathMatch(matcher Matcher) Matcher {
	return &urlPathMatcher{matcher: matcher}
}
