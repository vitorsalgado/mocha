package expect

// ToBe returns the result of the provided matcher.
func ToBe[V any](matcher Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "ToBe"
	m.Matches = matcher.Matches

	return m
}
