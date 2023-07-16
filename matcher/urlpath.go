package matcher

import (
	"fmt"
	"net/url"
	"strings"
)

type urlPathMatcher struct {
	matcher Matcher
}

func (m *urlPathMatcher) Match(v any) (Result, error) {
	var value string

	switch e := v.(type) {
	case *url.URL:
		value = e.Path
	case url.URL:
		value = e.Path
	case string:
		u, err := url.Parse(e)
		if err != nil {
			return Result{}, fmt.Errorf("urlpath: error parsing url: %s", err.Error())
		}

		value = u.Path
	case fmt.Stringer:
		u, err := url.Parse(e.String())
		if err != nil {
			return Result{}, fmt.Errorf("urlpath: error parsing url: %s", err.Error())
		}

		value = u.Path
	default:
		return Result{}, fmt.Errorf("urlpath: it only accepts the types: *url.URL | url.URL | string. got: %T", v)
	}

	res, err := m.matcher.Match(value)
	if err != nil {
		return Result{}, fmt.Errorf("urlpath: %w", err)
	}

	if res.Pass {
		return Result{Pass: true}, nil
	}

	return Result{Message: strings.Join([]string{"URLPath(", value, ") ", res.Message}, "")}, nil
}

// URLPath compares the URL path with the expected value and matches if they are equal.
// Comparison is case-insensitive.
func URLPath(expected string) Matcher {
	return URLPathMatch(GlobMatch(expected))
}

// URLPathf compares the URL path with the expected value and matches if they are equal.
// The expected value will be formatted according to the given format specifier.
// Comparison is case-insensitive.
func URLPathf(format string, a ...any) Matcher {
	return URLPathMatch(GlobMatch(fmt.Sprintf(format, a...)))
}

// URLPathMatch applies the provided matcher to the URL path.
func URLPathMatch(matcher Matcher) Matcher {
	return &urlPathMatcher{matcher: matcher}
}
