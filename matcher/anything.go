package matcher

func Anything[E any]() Matcher[E] {
	return func(v E, params Params) (bool, error) {
		return true, nil
	}
}
