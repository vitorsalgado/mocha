package matchers

// Is returns the result of the provided matcher.
func Is[V any](m Matcher[V]) Matcher[V] {
	return func(v V, params Args) (bool, error) {
		return m(v, params)
	}
}
