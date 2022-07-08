package to

// BeEmpty returns true if matcher value has zero length.
func BeEmpty[V any]() Matcher[V] {
	m := Matcher[V]{}
	m.Name = "IsEmpty"
	m.Matches = func(v V, args Args) (bool, error) {
		return HaveLen[V](0).Matches(v, args)
	}

	return m
}
