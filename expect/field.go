package expect

import "net/url"

// FormField matches a form field with the provided matcher.
func FormField(field string, matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "FormField"
	m.Matches = func(v any, args Args) (bool, error) {
		return matcher.Matches(v.(url.Values).Get(field), args)
	}

	return m
}
