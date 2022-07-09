package core

import (
	"sort"
	"sync"
)

type builtInStorage struct {
	data []*Mock
	mu   sync.Mutex
}

// NewStorage returns Mock storage implementation.
func NewStorage() Storage {
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
