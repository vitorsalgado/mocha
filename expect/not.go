package expect

// Not negates the provided matcher.
func Not[V any](matcher Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Not"
	m.Matches = func(v V, params Args) (bool, error) {
		result, err := matcher.Matches(v, params)
		return !result, err
	}

	return m
}
