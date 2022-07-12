package expect

// AnyOf matches when any of the given matchers returns true.
// Example:
//	AnyOf(EqualTo("test"),ToEqualFold("TEST"),ToContains("tes"))
func AnyOf(matchers ...Matcher) Matcher {
	m := Matcher{}
	m.Name = "AnyOf"
	m.Matches = func(v any, args Args) (bool, error) {
		for _, matcher := range matchers {
			if result, err := matcher.Matches(v, args); result || err != nil {
				return result, err
			}
		}

		return false, nil
	}

	return m
}
