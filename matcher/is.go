package matcher

func Is[V any](m Matcher[V]) Matcher[V] {
	return func(v V, params Params) (bool, error) {
		return m(v, params)
	}
}
