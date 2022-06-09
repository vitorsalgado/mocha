package arrays

func All[T any](arr []T, pred func(i T) bool) bool {
	for _, v := range arr {
		if !pred(v) {
			return false
		}
	}

	return true
}
