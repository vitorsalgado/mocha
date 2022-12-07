package mocha

import (
	"sort"
	"sync"
)

// mockStore is the definition for Mock repository.
type mockStore interface {
	// Save saves the Mock.
	Save(mock *Mock)

	// FetchEligible returns mocks that can be matched against requests.
	FetchEligible() []*Mock

	// FetchAll returns all stored Mock instances.
	FetchAll() []*Mock

	// Delete removes a Mock by its ID.
	Delete(id int)

	// DeleteBySource removes mocks by its source.
	DeleteBySource(source string)

	// Flush removes all stored mocks.
	Flush()
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

func (repo *builtInStore) FetchEligible() []*Mock {
	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (repo *builtInStore) FetchAll() []*Mock {
	return repo.data
}

func (repo *builtInStore) Delete(id int) {
	index := -1
	for i, m := range repo.data {
		if m.ID == id {
			index = i
			break
		}
	}

	repo.data = repo.data[:index+copy(repo.data[index:], repo.data[index+1:])]
}

func (repo *builtInStore) DeleteBySource(source string) {
	index := -1
	for i, m := range repo.data {
		if m.Source == source {
			index = i
			break
		}
	}

	repo.data = repo.data[:index+copy(repo.data[index:], repo.data[index+1:])]
}

func (repo *builtInStore) Flush() {
	repo.data = nil
	repo.data = make([]*Mock, 0)
}
