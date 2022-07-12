package expect

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek(matcher Matcher, action func(v any) error) Matcher {
	m := Matcher{}
	m.Name = "Peek"
	m.Matches = func(v any, params Args) (bool, error) {
		err := action(v)
		if err != nil {
			return false, err
		}

		return matcher.Matches(v, params)
	}

	return m
}
