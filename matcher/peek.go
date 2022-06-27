package matcher

func Peek[V any](m Matcher[V], action func(v V) error) Matcher[V] {
	return func(v V, params Params) (bool, error) {
		err := action(v)
		if err != nil {
			return false, err
		}

		return m(v, params)
	}
}
