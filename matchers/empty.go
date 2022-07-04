package matchers

// IsEmpty returns true if matcher value has zero length.
func IsEmpty[V any]() Matcher[V] {
	m := Matcher[V]{}
	m.Name = "IsEmpty"
	m.Matches = func(v V, args Args) (bool, error) {
		return Len[V](0).Matches(v, args)
	}

	return m
}
