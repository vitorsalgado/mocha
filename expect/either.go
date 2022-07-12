package expect

// EitherMatcherBuilder is a builder for Either matcher.
// Prefer to use the Either() function.
type EitherMatcherBuilder struct {
	first Matcher
}

// Either matches true when any of the two given matchers returns true.
func Either(first Matcher) *EitherMatcherBuilder {
	return &EitherMatcherBuilder{first}
}

// Or sets the second matcher
func (e *EitherMatcherBuilder) Or(second Matcher) Matcher {
	m := Matcher{}
	m.Name = "Either"
	m.Matches = func(v any, args Args) (bool, error) {
		r1, err := e.first.Matches(v, args)
		if err != nil {
			return false, err
		}

		r2, err := second.Matches(v, args)

		return r1 || r2, err
	}

	return m
}
