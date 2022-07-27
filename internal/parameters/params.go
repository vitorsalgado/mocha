// Package parameters implements a simple in-memory key/value store, used internally by Mocha.
package parameters

type (
	// Params defines the contract to an in-memory storage for generic parameters.
	Params interface {
		// Get returns the parameter by its key.
		Get(key string) (any, bool)

		// GetAll returns all stored parameters.
		GetAll() map[string]any

		// Set sets a parameter.
		Set(key string, dep any)

		// Remove removes a parameter by its key.
		Remove(key string)

		// Has checks if a parameter with the given key exists.
		Has(key string) bool
	}

	inMemoryParams struct {
		data map[string]any
	}
)

// New returns a Params concrete implementation.
func New() Params {
	return &inMemoryParams{data: make(map[string]any)}
}

func (p inMemoryParams) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p inMemoryParams) GetAll() map[string]any {
	return p.data
}

func (p inMemoryParams) Set(key string, dep any) {
	p.data[key] = dep
}

func (p inMemoryParams) Remove(key string) {
	delete(p.data, key)
}

func (p inMemoryParams) Has(key string) bool {
	_, ok := p.data[key]

	return ok
}
