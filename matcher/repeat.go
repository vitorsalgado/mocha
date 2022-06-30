package matcher

func Repeat[V any](times int) Matcher[V] {
	count := 0

	return func(_ V, params Args) (bool, error) {
		count++

		return count <= times, nil
	}
}
