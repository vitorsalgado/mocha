package matchers

// Either matches true when any of the two given matchers returns true.
func Either[V any](first Matcher[V], second Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Either"
	m.Matches = func(v V, args Args) (bool, error) {
		r1, err := first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := second.Matches(v, args)

		return r1 || r2, err
	}

	return m
}
