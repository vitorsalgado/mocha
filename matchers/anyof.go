package matchers

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),EqualFold("TEST"),Contains("tes"))
func AnyOf[E any](matchers ...Matcher[E]) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		for _, m := range matchers {
			if result, err := m(v, params); result || err != nil {
				return result, err
			}
		}

		return false, nil
	}
}
