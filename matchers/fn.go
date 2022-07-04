package matchers

// Fn creates an anonymous Matcher using the given function.
func Fn[V any](fn func(v V, a Args) (bool, error)) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Anonymous"
	m.Matches = fn

	return m
}
