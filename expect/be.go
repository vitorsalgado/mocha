package expect

// ToBe returns the result of the provided matcher.
func ToBe(matcher Matcher) Matcher {
	m := Matcher{}
	m.Name = "ToBe"
	m.Matches = matcher.Matches

	return m
}
