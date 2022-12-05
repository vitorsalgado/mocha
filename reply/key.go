package reply

// K is a http.Request context key.
type K int8

// http.Request context keys used by reply implementations.
const (
	// KArg is a key to retrieve Arg.
	KArg K = iota
)
