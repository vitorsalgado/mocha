package expect

// AllOf matches when all the given matchers returns true.
// Example:
//	AllOf(EqualTo("test"),ToEqualFold("test"),ToContains("tes"))
func AllOf(matchers ...Matcher) Matcher {
	m := Matcher{}
	m.Name = "AllOf"
	m.Matches = func(v any, args Args) (bool, error) {
		for _, matcher := range matchers {
			if result, err := matcher.Matches(v, args); !result || err != nil {
				return result, err
			}
		}

		return true, nil
	}

	return m
}
