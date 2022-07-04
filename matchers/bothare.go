package matchers

// BothAreMatcher allows a more fluent definition of BothAre matcher.
type BothAreMatcher[E any] struct {
	first Matcher[E]
}

// BothAre matches true when both given matchers evaluates to true.
func BothAre[E any](first Matcher[E]) BothAreMatcher[E] {
	return BothAreMatcher[E]{first: first}
}

// And sets the second matcher.
func (ba BothAreMatcher[E]) And(second Matcher[E]) Matcher[E] {
	return func(v E, params Args) (bool, error) {
		r1, err := ba.first(v, params)
		if err != nil {
			return false, err
		}

		r2, err := second(v, params)

		return r1 && r2, err
	}
}