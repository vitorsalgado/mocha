package mocha

import (
	"sync"
)

var _ Params = (*paramsStore)(nil)

// Params defines a contract for a generic parameters repository.
type Params interface {
	// Get returns the parameter by its key.
	Get(k string) (datum any, err error)

	// GetAll returns all stored parameters.
	GetAll() (map[string]any, error)

	// Set sets a parameter.
	Set(k string, v any) error

	MustSet(k string, v any)

	// Remove removes a parameter by its key.
	Remove(k string) error

	// Has checks if a parameter with the given key exists.
	Has(k string) (bool, error)
}

type paramsStore struct {
	data map[string]any
	mu   sync.RWMutex
}

func newInMemoryParameters() Params {
	return &paramsStore{data: make(map[string]any)}
}

func (p *paramsStore) Get(key string) (datum any, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	datum = p.data[key]

	return
}

func (p *paramsStore) GetAll() (map[string]any, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.data, nil
}

func (p *paramsStore) Set(key string, dep any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data[key] = dep
	return nil
}

func (p *paramsStore) MustSet(key string, dep any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.data[key] = dep
}

func (p *paramsStore) Remove(key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.data, key)
	return nil
}

func (p *paramsStore) Has(key string) (bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.data[key]
	return ok, nil
}
