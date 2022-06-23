package mocha

import (
	"errors"
	"github.com/vitorsalgado/mocha/mock"
)

type Scoped struct {
	store mock.Storage
	mocks []int32
}

var (
	ErrScopeNotDone = errors.New("there are still mocks that were not called")
)

func NewScoped(repo mock.Storage, mocks []int32) *Scoped {
	return &Scoped{store: repo, mocks: mocks}
}

func (s *Scoped) IsDone() bool {
	for _, key := range s.mocks {
		m := s.store.GetByID(key)
		if !m.Called() {
			return false
		}
	}

	return true
}

func (s *Scoped) Pending() []mock.Mock {
	return s.store.Pending(s.mocks)
}

func (s *Scoped) IsPending() bool {
	return len(s.store.Pending(s.mocks)) > 0
}

func (s *Scoped) Clean() {
	for _, key := range s.mocks {
		s.store.Delete(key)
	}

	s.mocks = make([]int32, 0)
}

func (s *Scoped) Done() error {
	if s.IsPending() {
		return ErrScopeNotDone
	}

	return nil
}
