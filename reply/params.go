package reply

// Params defines a contract for a generic parameters repository.
type Params interface {
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

type paramsStore struct {
	data map[string]any
}

// Parameters returns a Params concrete implementation.
func Parameters() Params {
	return &paramsStore{data: make(map[string]any)}
}

func (p paramsStore) Get(key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p paramsStore) GetAll() map[string]any {
	return p.data
}

func (p paramsStore) Set(key string, dep any) {
	p.data[key] = dep
}

func (p paramsStore) Remove(key string) {
	delete(p.data, key)
}

func (p paramsStore) Has(key string) bool {
	_, ok := p.data[key]

	return ok
}
