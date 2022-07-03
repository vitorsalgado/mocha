package matcher

// Not negates the provided matcher.
func Not[V any](m Matcher[V]) Matcher[V] {
	return func(v V, params Args) (bool, error) {
		result, err := m(v, params)
		return !result, err
	}
}
