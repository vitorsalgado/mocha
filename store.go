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

func (s *builtInStore) Save(mock *Mock) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = append(s.data, mock)

	sort.SliceStable(s.data, func(a, b int) bool {
		return s.data[a].Priority < s.data[b].Priority
	})
}

func (s *builtInStore) Get(id string) *Mock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	for _, datum := range s.data {
		if datum.ID == id {
			return datum
		}
	}

	return nil
}

func (s *builtInStore) GetEligible() []*Mock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	mocks := make([]*Mock, 0, len(s.data))

	for _, mock := range s.data {
		if mock.Enabled {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (s *builtInStore) GetAll() []*Mock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	return s.data
}

func (s *builtInStore) Delete(id string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	index := -1
	for i, m := range s.data {
		if m.ID == id {
			index = i
			break
		}
	}

	s.data = s.data[:index+copy(s.data[index:], s.data[index+1:])]
}

func (s *builtInStore) DeleteExternal() {
	ids := make([]string, 0, len(s.data))

	for _, m := range s.data {
		if len(m.Source) > 0 {
			ids = append(ids, m.ID)
		}
	}

	for _, id := range ids {
		s.Delete(id)
	}
}

func (s *builtInStore) DeleteAll() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = nil
	s.data = make([]*Mock, 0)
}
