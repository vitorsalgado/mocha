package matchers

// BothAreMatcher allows a more fluent definition of BothAre matcher.
type BothAreMatcher[V any] struct {
	Matcher[V]
	first  Matcher[V]
	second Matcher[V]
}

// BothAre matches true when both given matchers evaluates to true.
func BothAre[V any](first Matcher[V]) *BothAreMatcher[V] {
	m := &BothAreMatcher[V]{first: first}
	m.Name = "BothAre"
	m.Matches = func(v V, args Args) (bool, error) {
		r1, err := m.first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := m.second.Matches(v, args)

		return r1 && r2, err
	}

	return m
}

// And sets the second matcher.
func (ba *BothAreMatcher[E]) And(second Matcher[E]) BothAreMatcher[E] {
	ba.second = second
	return *ba
}
