package expect

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf[V any](matchers ...Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "AnyOf"
	m.Matches = func(v V, args Args) (bool, error) {
		for _, matcher := range matchers {
			if result, err := matcher.Matches(v, args); result || err != nil {
				return result, err
			}
		}

		return false, nil
	}

	return m
}
