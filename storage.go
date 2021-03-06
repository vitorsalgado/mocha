package mocha

import (
	"sort"
	"sync"
)

// storage is the definition for Mock repository.
type storage interface {
	// Save saves the Mock.
	Save(mock *Mock)

	// FetchEligible returns mocks that can be matched against requests.
	FetchEligible() []*Mock

	// FetchAll returns all stored Mock instances.
	FetchAll() []*Mock

	// Delete removes a Mock by its ID.
	Delete(id int)

	// Flush removes all stored mocks.
	Flush()
}

type builtInStorage struct {
	data []*Mock
	mu   sync.Mutex
}

// newStorage returns Mock storage implementation.
func newStorage() storage {
	return &builtInStorage{data: make([]*Mock, 0)}
}

func (repo *builtInStorage) Save(mock *Mock) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data = append(repo.data, mock)

	sort.SliceStable(repo.data, func(a, b int) bool {
		return repo.data[a].Priority < repo.data[b].Priority
	})
}

func (repo *builtInStorage) FetchEligible() []*Mock {
	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (repo *builtInStorage) FetchAll() []*Mock {
	return repo.data
}

func (repo *builtInStorage) Delete(id int) {
	index := -1
	for i, m := range repo.data {
		if m.ID == id {
			index = i
			break
		}
	}

	repo.data = repo.data[:index+copy(repo.data[index:], repo.data[index+1:])]
}

func (repo *builtInStorage) Flush() {
	repo.data = nil
	repo.data = make([]*Mock, 0)
}
