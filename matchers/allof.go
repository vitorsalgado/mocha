package matchers

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),EqualFold("test"),Contains("tes"))
func AllOf[V any](matchers ...Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "AllOf"
	m.Matches = func(v V, args Args) (bool, error) {
		for _, matcher := range matchers {
			if result, err := matcher.Matches(v, args); !result || err != nil {
				return result, err
			}
		}

		return true, nil
	}

	return m
}
