package matchers

// Anything returns true all the time.
func Anything[V any]() Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Anything"
	m.Matches = func(_ V, _ Args) (bool, error) {
		return true, nil
	}

	return m
}
