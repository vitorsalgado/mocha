package lib

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

	// MustSet sets a parameter and if any error occurs, it should be discarded or panic.
	MustSet(k string, v any)

	// Remove removes a parameter by its key.
	Remove(k string) error

	// Has checks if a parameter with the given key exists.
	Has(k string) (bool, error)
}

type paramsStore struct {
	data    map[string]any
	rwMutex sync.RWMutex
}

func NewInMemoryParameters() Params {
	return &paramsStore{data: make(map[string]any)}
}

func (p *paramsStore) Get(key string) (datum any, err error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	datum = p.data[key]

	return
}

func (p *paramsStore) GetAll() (map[string]any, error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	return p.data, nil
}

func (p *paramsStore) Set(key string, dep any) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.data[key] = dep

	return nil
}

func (p *paramsStore) MustSet(key string, dep any) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	p.data[key] = dep
}

func (p *paramsStore) Remove(key string) error {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	delete(p.data, key)

	return nil
}

func (p *paramsStore) Has(key string) (bool, error) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	_, ok := p.data[key]

	return ok, nil
}
