package dzstd

import (
	"context"
	"sort"
	"sync"
)

var _ MockRepository[*BaseMock] = (*InMemMockStore[*BaseMock])(nil)

type MockRepository[TMock Mock] interface {
	FindEligible(ctx context.Context, done chan struct{}) (<-chan TMock, error)

	Save(ctx context.Context, mock ...TMock) error

	Get(ctx context.Context, id string) (TMock, error)
	GetAll(ctx context.Context) ([]TMock, error)
	GetByIDs(ctx context.Context, ids ...string) ([]TMock, error)

	Delete(ctx context.Context, ids ...string) error
	DeleteExternal(ctx context.Context) error
	DeleteAll(ctx context.Context) error
}

type InMemMockStore[TMock Mock] struct {
	data    []TMock
	rwMutex sync.RWMutex
}

// NewStore returns Mock InMemMockStore implementation.
func NewStore[TMock Mock]() *InMemMockStore[TMock] {
	return &InMemMockStore[TMock]{data: make([]TMock, 0)}
}

func (s *InMemMockStore[TMock]) FindEligible(_ context.Context, done chan struct{}) (<-chan TMock, error) {
	s.rwMutex.RLock()
	mocks := make([]TMock, len(s.data))
	copy(mocks, s.data)
	s.rwMutex.RUnlock()

	out := make(chan TMock)
	go func(mocks []TMock) {
		defer close(out)
		for _, mock := range mocks {
			if mock.IsEnabled() {
				select {
				case <-done:
					return
				case out <- mock:
				}
			}
		}
	}(mocks)

	return out, nil
}

func (s *InMemMockStore[TMock]) Save(_ context.Context, mocks ...TMock) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = append(s.data, mocks...)

	sort.SliceStable(s.data, func(a, b int) bool {
		return s.data[a].GetPriority() < s.data[b].GetPriority()
	})

	return nil
}

func (s *InMemMockStore[TMock]) Get(_ context.Context, id string) (TMock, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	for _, datum := range s.data {
		if datum.GetID() == id {
			return datum, nil
		}
	}

	var result TMock
	return result, nil
}

func (s *InMemMockStore[TMock]) GetByIDs(ctx context.Context, ids ...string) ([]TMock, error) {
	mocks := make([]TMock, 0, len(ids))
	for _, id := range ids {
		m, err := s.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		mocks = append(mocks, m)
	}

	return mocks, nil
}

func (s *InMemMockStore[TMock]) GetAll(_ context.Context) ([]TMock, error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	ret := make([]TMock, len(s.data))
	copy(ret, s.data)

	return ret, nil
}

func (s *InMemMockStore[TMock]) Delete(_ context.Context, ids ...string) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for _, id := range ids {
		index := -1
		for i, m := range s.data {
			if m.GetID() == id {
				index = i
				break
			}
		}

		s.data = s.data[:index+copy(s.data[index:], s.data[index+1:])]
	}

	return nil
}

func (s *InMemMockStore[TMock]) DeleteExternal(ctx context.Context) error {
	ids := make([]string, 0, len(s.data))

	for _, m := range s.data {
		if len(m.GetSource()) > 0 {
			ids = append(ids, m.GetID())
		}
	}

	for _, id := range ids {
		if err := s.Delete(ctx, id); err != nil {
			return err
		}
	}

	return nil
}

func (s *InMemMockStore[TMock]) DeleteAll(_ context.Context) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = nil
	s.data = make([]TMock, 0)

	return nil
}
