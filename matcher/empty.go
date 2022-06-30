package matcher

func IsEmpty[V any]() Matcher[V] {
	return func(v V, params Args) (bool, error) {
		return Len[V](0)(v, params)
	}
}
