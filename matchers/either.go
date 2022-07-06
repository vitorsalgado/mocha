package matchers

// EitherMatcherBuilder is a builder for Either matcher.
// Prefer to use the Either() function.
type EitherMatcherBuilder[V any] struct {
	first Matcher[V]
}

// Either matches true when any of the two given matchers returns true.
func Either[V any](first Matcher[V]) *EitherMatcherBuilder[V] {
	return &EitherMatcherBuilder[V]{first}
}

// Or sets the second matcher
func (e *EitherMatcherBuilder[V]) Or(second Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Either"
	m.Matches = func(v V, args Args) (bool, error) {
		r1, err := e.first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := second.Matches(v, args)

		return r1 || r2, err
	}

	return m
}
