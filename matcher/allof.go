package matcher

// AllOf matches when all of the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),EqualFold("test"),Contains("tes"))
func AllOf[E any](matchers ...Matcher[E]) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		for _, m := range matchers {
			if result, err := m(v, params); !result || err != nil {
				return result, err
			}
		}

		return true, nil
	}
}
