package types

type RawValueTypes interface {
	map[string]any | []any | string
}

// RawValue describes a Matcher in its raw format.
type RawValue []any

func (r RawValue) Arguments() []any {
	if len(r) > 1 {
		return r[1:]
	}

	return nil
}

// Persist represents a mock server object that can be saved and loaded.
type Persist interface {
	// Raw describes an object in a way that it can be saved and loaded.
	Raw() RawValue
}
