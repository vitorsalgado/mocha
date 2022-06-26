package matcher

type bothAre[E any] struct {
	first Matcher[E]
}

func BothAre[E any](first Matcher[E]) bothAre[E] {
	return bothAre[E]{first: first}
}

func (ba bothAre[E]) And(second Matcher[E]) Matcher[E] {
	return func(v E, params Params) (bool, error) {
		r1, err := ba.first(v, params)
		if err != nil {
			return false, err
		}

		r2, err := second(v, params)

		return r1 && r2, err
	}
}
