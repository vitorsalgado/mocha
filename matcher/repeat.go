package matcher

func Repeat(times int) Matcher[any] {
	count := 0

	return func(v any, params Params) (bool, error) {
		count++

		return count <= times, nil
	}
}
