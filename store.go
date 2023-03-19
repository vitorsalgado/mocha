package mocha

import (
	"sort"
	"sync"
)

var _ mockStore = (*builtInStore)(nil)

// mockStore is the definition of a Mock repository.
type mockStore interface {
	// Save saves the Mock.
	Save(mock *Mock)

	// Get returns a Mock by ID.
	Get(id string) *Mock

	// GetEligible returns mocks that are eligible to be matched and served.
	GetEligible() []*Mock

	// GetAll returns all stored Mock instances.
	GetAll() []*Mock

	// Delete removes a Mock by its ID.
	Delete(id string)

	// DeleteExternal removes mocks set by external components, like Loader.
	// Mostly used internally.
	DeleteExternal()

	// DeleteAll removes all stored mocks.
	DeleteAll()
}

type builtInStore struct {
	data    []*Mock
	rwMutex sync.RWMutex
}

// newStore returns Mock mockStore implementation.
func newStore() mockStore {
	return &builtInStore{data: make([]*Mock, 0)}
}

func (repo *builtInStore) Save(mock *Mock) {
	repo.rwMutex.Lock()
	defer repo.rwMutex.Unlock()

	repo.data = append(repo.data, mock)

	sort.SliceStable(repo.data, func(a, b int) bool {
		return repo.data[a].Priority < repo.data[b].Priority
	})
}

func (repo *builtInStore) Get(id string) *Mock {
	repo.rwMutex.RLock()
	defer repo.rwMutex.RUnlock()

	for _, datum := range repo.data {
		if datum.ID == id {
			return datum
		}
	}

	return nil
}

func (repo *builtInStore) GetEligible() []*Mock {
	repo.rwMutex.RLock()
	defer repo.rwMutex.RUnlock()

	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (repo *builtInStore) GetAll() []*Mock {
	repo.rwMutex.RLock()
	defer repo.rwMutex.RUnlock()

	return repo.data
}

func (repo *builtInStore) Delete(id string) {
	repo.rwMutex.Lock()
	defer repo.rwMutex.Unlock()

	index := -1
	for i, m := range repo.data {
		if m.ID == id {
			index = i
			break
		}
	}

	repo.data = repo.data[:index+copy(repo.data[index:], repo.data[index+1:])]
}

func (repo *builtInStore) DeleteExternal() {
	repo.rwMutex.Lock()
	defer repo.rwMutex.Unlock()

	data := make([]*Mock, 0)

	for _, m := range repo.data {
		if m.Source == "" {
			data = append(data, m)
		}
	}

	repo.data = data
}

func (repo *builtInStore) DeleteAll() {
	repo.rwMutex.Lock()
	defer repo.rwMutex.Unlock()

	repo.data = nil
	repo.data = make([]*Mock, 0)
}
