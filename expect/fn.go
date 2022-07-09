package expect

// Func creates an anonymous Matcher using the given function.
func Func[V any](fn func(v V, a Args) (bool, error)) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Func"
	m.Matches = fn

	return m
}
