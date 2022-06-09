package mocha

import "sort"

type MockRepository interface {
	Save(mock Mock)
	GetAllSorted() []*Mock
	GetByID(id int32) *Mock
}

type InMemoryMockRepository struct {
	data map[int32]*Mock
}

func NewMockRepository() MockRepository {
	return &InMemoryMockRepository{data: make(map[int32]*Mock)}
}

func (repo *InMemoryMockRepository) Save(mock Mock) {
	repo.data[mock.ID] = &mock
}

func (repo *InMemoryMockRepository) GetByID(id int32) *Mock {
	return repo.data[id]
}

func (repo *InMemoryMockRepository) GetAllSorted() []*Mock {
	var mocks = make([]*Mock, 0, len(repo.data))

	for _, mock := range repo.data {
		mocks = append(mocks, mock)
	}

	sort.SliceStable(mocks, func(a, b int) bool {
		return mocks[a].Priority < mocks[b].Priority
	})

	return mocks
}
