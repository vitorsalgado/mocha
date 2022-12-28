package mocha

import "context"

// Params defines a contract for a generic parameters repository.
type Params interface {
	// Get returns the parameter by its key.
	Get(ctx context.Context, k string) (any, bool)

	// GetAll returns all stored parameters.
	GetAll(ctx context.Context) map[string]any

	// Set sets a parameter.
	Set(ctx context.Context, k string, v any)

	// Remove removes a parameter by its key.
	Remove(ctx context.Context, k string)

	// Has checks if a parameter with the given key exists.
	Has(ctx context.Context, k string) bool
}

type paramsStore struct {
	data map[string]any
}

// Parameters returns a Params concrete implementation.
func Parameters() Params {
	return &paramsStore{data: make(map[string]any)}
}

func (p paramsStore) Get(_ context.Context, key string) (any, bool) {
	val, ok := p.data[key]
	return val, ok
}

func (p paramsStore) GetAll(_ context.Context) map[string]any {
	return p.data
}

func (p paramsStore) Set(_ context.Context, key string, dep any) {
	p.data[key] = dep
}

func (p paramsStore) Remove(_ context.Context, key string) {
	delete(p.data, key)
}

func (p paramsStore) Has(_ context.Context, key string) bool {
	_, ok := p.data[key]

	return ok
}
