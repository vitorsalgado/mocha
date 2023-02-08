package mocha

import (
	"context"
	"sync"
)

// Params defines a contract for a generic parameters repository.
type Params interface {
	// Get returns the parameter by its key.
	Get(ctx context.Context, k string) (datum any, exists bool, err error)

	// GetAll returns all stored parameters.
	GetAll(ctx context.Context) (map[string]any, error)

	// Set sets a parameter.
	Set(ctx context.Context, k string, v any) error

	// Remove removes a parameter by its key.
	Remove(ctx context.Context, k string) error

	// Has checks if a parameter with the given key exists.
	Has(ctx context.Context, k string) (bool, error)
}

type paramsStore struct {
	data map[string]any
	mu   sync.RWMutex
}

func newInMemoryParameters() Params {
	return &paramsStore{data: make(map[string]any)}
}

func (p *paramsStore) Get(_ context.Context, key string) (datum any, exists bool, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	datum, exists = p.data[key]
	return
}

func (p *paramsStore) GetAll(_ context.Context) (map[string]any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.data, nil
}

func (p *paramsStore) Set(_ context.Context, key string, dep any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data[key] = dep
	return nil
}

func (p *paramsStore) Remove(_ context.Context, key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.data, key)
	return nil
}

func (p *paramsStore) Has(_ context.Context, key string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.data[key]
	return ok, nil
}
