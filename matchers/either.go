package matchers

// Either matches true when any of the two given matchers returns true.
func Either[V any](first Matcher[V], second Matcher[V]) Matcher[V] {
	return func(v V, params Args) (bool, error) {
		r1, err := first(v, params)
		if err != nil {
			return false, err
		}

		r2, err := second(v, params)

		return r1 || r2, err
	}
}
