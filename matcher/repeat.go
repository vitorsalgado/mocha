package matcher

func Repeat[V any](times int) Matcher[V] {
	count := 0

	return func(_ V, params Params) (bool, error) {
		count++

		return count <= times, nil
	}
}
