package mock

import (
	"sort"
	"sync"
)

type (
	BuiltInStorage struct {
		data []*Mock
		mu   sync.Mutex
	}
)

func NewStorage() Storage {
	return &BuiltInStorage{data: make([]*Mock, 0)}
}

func (repo *BuiltInStorage) Save(mock *Mock) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data = append(repo.data, mock)

	sort.SliceStable(repo.data, func(a, b int) bool {
		return repo.data[a].Priority < repo.data[b].Priority
	})
}

func (repo *BuiltInStorage) FetchEligible() []*Mock {
	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (repo *BuiltInStorage) FetchAll() []*Mock {
	return repo.data
}

func (repo *BuiltInStorage) Delete(id int) {
	index := -1
	for i, m := range repo.data {
		if m.ID == id {
			index = i
			break
		}
	}

	repo.data = repo.data[:index+copy(repo.data[index:], repo.data[index+1:])]
}

func (repo *BuiltInStorage) Flush() {
	repo.data = nil
	repo.data = make([]*Mock, 0)
}
