package mock

import (
	"sort"
	"sync"
)

type (
	BuiltInStorage struct {
		data map[int]*Mock
		mu   sync.Mutex
	}
)

func NewStorage() Storage {
	return &BuiltInStorage{data: make(map[int]*Mock)}
}

func (repo *BuiltInStorage) Save(mock *Mock) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.data[mock.ID] = mock
}

func (repo *BuiltInStorage) FetchByID(id int) *Mock {
	return repo.data[id]
}

func (repo *BuiltInStorage) FetchByIDs(ids ...int) []*Mock {
	size := len(ids)
	arr := make([]*Mock, size, size)

	for i, id := range ids {
		arr[i] = repo.FetchByID(id)
	}

	return arr
}

func (repo *BuiltInStorage) FetchEligibleSorted() []*Mock {
	mocks := make([]*Mock, 0)

	for _, mock := range repo.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	sort.SliceStable(mocks, func(a, b int) bool {
		return mocks[a].Priority < mocks[b].Priority
	})

	return mocks
}

func (repo *BuiltInStorage) FetchAll() []*Mock {
	size := len(repo.data)
	mocks := make([]*Mock, size, size)
	i := 0

	for _, m := range repo.data {
		mocks[i] = m
		i++
	}

	return mocks
}

func (repo *BuiltInStorage) Delete(id int) {
	delete(repo.data, id)
}

func (repo *BuiltInStorage) Flush() {
	for key := range repo.data {
		delete(repo.data, key)
	}
}
