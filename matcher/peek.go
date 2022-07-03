package matcher

// Peek will return the result of the given matcher, after executing the provided function.
// Peek can be used to check the matcher argument.
func Peek[V any](m Matcher[V], action func(v V) error) Matcher[V] {
	return func(v V, params Args) (bool, error) {
		err := action(v)
		if err != nil {
			return false, err
		}

		return m(v, params)
	}
}
