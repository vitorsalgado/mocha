package expect

// Func creates an anonymous Matcher using the given function.
func Func(fn func(v any, a Args) (bool, error)) Matcher {
	m := Matcher{}
	m.Name = "Func"
	m.Matches = fn

	return m
}
