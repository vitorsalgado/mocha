package matcher

func AnyOf[E any](matchers ...Matcher[E]) Matcher[E] {
	return func(v E, params Params) (bool, error) {
		for _, m := range matchers {
			if result, err := m(v, params); result || err != nil {
				return result, err
			}
		}

		return false, nil
	}
}
