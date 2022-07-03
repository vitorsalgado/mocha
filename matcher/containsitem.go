package matcher

// ContainsItem returns when the given value is contained in the matcher slice.
func ContainsItem[V any](value V) Matcher[[]V] {
	return func(v []V, params Args) (bool, error) {
		for _, entry := range v {
			if r, err := EqualTo(value)(entry, params); r || err != nil {
				return r, err
			}
		}

		return false, nil
	}
}
