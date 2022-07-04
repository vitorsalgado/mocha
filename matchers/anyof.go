package matchers

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),EqualFold("TEST"),Contains("tes"))
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
