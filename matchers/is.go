package matchers

// Is returns the result of the provided matcher.
func Is[V any](matcher Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Is"
	m.Matches = matcher.Matches

	return m
}
