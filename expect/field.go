package expect

import "net/url"

// FormField matches a form field with the provided matcher.
func FormField(field string, matcher Matcher[string]) Matcher[url.Values] {
	m := Matcher[url.Values]{}
	m.Name = "FormField"
	m.Matches = func(v url.Values, args Args) (bool, error) {
		return matcher.Matches(v.Get(field), args)
	}

	return m
}
