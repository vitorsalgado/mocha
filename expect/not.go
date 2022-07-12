package expect

// Not negates the provided matcher.
func Not(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "Not"
	m.Matches = func(v any, params Args) (bool, error) {
		result, err := matcher.Matches(v, params)
		return !result, err
	}

	return m
}
