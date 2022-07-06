package matchers

// BothMatcherBuilder is a builder for Both matcher.
// Use .Both() function to create a new Both matcher.
type BothMatcherBuilder[V any] struct {
	first Matcher[V]
}

// Both matches true when both given matchers evaluates to true.
func Both[V any](first Matcher[V]) *BothMatcherBuilder[V] {
	m := &BothMatcherBuilder[V]{first: first}

	return m
}

// And sets the second matcher.
func (ba *BothMatcherBuilder[E]) And(second Matcher[E]) Matcher[E] {
	m := Matcher[E]{}
	m.Name = "Both"
	m.Matches = func(v E, args Args) (bool, error) {
		r1, err := ba.first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := second.Matches(v, args)

		return r1 && r2, err
	}

	return m
}
