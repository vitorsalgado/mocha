package expect

// ToBeEmpty returns true if matcher value has zero length.
func ToBeEmpty() Matcher {
	m := Matcher{}
	m.Name = "Empty"
	m.Matches = func(v any, args Args) (bool, error) {
		return ToHaveLen(0).Matches(v, args)
	}

	return m
}
