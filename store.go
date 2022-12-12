package mocha

import (
	"sort"
	"sync"
)

var _ mockStore = (*builtInStore)(nil)

// mockStore is the definition for Mock repository.
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
	data []*Mock
	mu   sync.Mutex
}

// newStore returns Mock mockStore implementation.
func newStore() mockStore {
	return &builtInStore{data: make([]*Mock, 0)}
}

func (repo *builtInStore) Save(mock *Mock) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data = append(repo.data, mock)

	sort.SliceStable(repo.data, func(a, b int) bool {
		return repo.data[a].Priority < repo.data[b].Priority
	})
}

func (repo *builtInStore) Get(id string) *Mock {
	for _, datum := range repo.data {
		if datum.ID == id {
			return datum
		}
	}

	return nil
}

func (repo *builtInStore) GetEligible() []*Mock {
	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (repo *builtInStore) GetAll() []*Mock {
	return repo.data
}

func (repo *builtInStore) Delete(id string) {
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
	for i, m := range repo.data {
		if m.Source != "" {
			repo.data = repo.data[:i+copy(repo.data[i:], repo.data[i+1:])]
		}
	}
}

func (repo *builtInStore) DeleteAll() {
	repo.data = nil
	repo.data = make([]*Mock, 0)
}
