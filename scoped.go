package mocha

import (
	"errors"
	"github.com/vitorsalgado/mocha/mock"
)

type Scoped struct {
	storage mock.Storage
	mocks   []int32
}

var (
	ErrScopeNotDone = errors.New("there are still mocks that were not called")
)

func Scope(repo mock.Storage, mocks []int32) *Scoped {
	return &Scoped{storage: repo, mocks: mocks}
}

func (s *Scoped) IsDone() bool {
	for _, key := range s.mocks {
		m := s.storage.FetchByID(key)
		if !m.Called() {
			return false
		}
	}

	return true
}

func (s *Scoped) Pending() []mock.Mock {
	return s.storage.Pending(s.mocks)
}

func (s *Scoped) IsPending() bool {
	return len(s.storage.Pending(s.mocks)) > 0
}

func (s *Scoped) Clean() {
	for _, key := range s.mocks {
		s.storage.Delete(key)
	}

	s.mocks = make([]int32, 0)
}

func (s *Scoped) Done() error {
	if s.IsPending() {
		return ErrScopeNotDone
	}

	return nil
}
