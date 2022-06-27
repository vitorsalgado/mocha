package matcher

func Either[V any](first Matcher[V], second Matcher[V]) Matcher[V] {
	return func(v V, params Params) (bool, error) {
		r1, err := first(v, params)
		if err != nil {
			return false, err
		}

		r2, err := second(v, params)

		return r1 || r2, err
	}
}
