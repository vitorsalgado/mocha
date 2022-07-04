package matchers

// ContainsItem returns when the given value is contained in the matcher slice.
func ContainsItem[V any](value V) Matcher[[]V] {
	m := Matcher[[]V]{}
	m.Name = "ContainsItem"
	m.Matches = func(v []V, args Args) (bool, error) {
		for _, entry := range v {
			if r, err := EqualTo(value).Matches(entry, emptyArgs()); r || err != nil {
				return r, err
			}
		}

		return false, nil
	}

	return m
}
