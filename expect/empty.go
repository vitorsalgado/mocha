package expect

// ToBeEmpty returns true if matcher value has zero length.
func ToBeEmpty[V any]() Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Empty"
	m.Matches = func(v V, args Args) (bool, error) {
		return ToHaveLen[V](0).Matches(v, args)
	}

	return m
}
