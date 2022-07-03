package matcher

// Anything returs true all the time.
func Anything[E any]() Matcher[E] {
	return func(v E, params Args) (bool, error) {
		return true, nil
	}
}
