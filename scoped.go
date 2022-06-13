package mocha

import (
	"errors"
)

type Scoped struct {
	repo  MockRepository
	mocks []int32
}

var (
	ErrScopeNotDone = errors.New("there are still mocks that were not called")
)

func NewScoped(repo MockRepository, mocks []int32) *Scoped {
	return &Scoped{repo: repo, mocks: mocks}
}

func (s *Scoped) IsDone() bool {
	for _, key := range s.mocks {
		m := s.repo.GetByID(key)
		if !m.Called() {
			return false
		}
	}

	return true
}

func (s *Scoped) Pending() []Mock {
	return s.repo.Pending(s.mocks)
}

func (s *Scoped) IsPending() bool {
	return len(s.repo.Pending(s.mocks)) > 0
}

func (s *Scoped) Clean() {
	for _, key := range s.mocks {
		s.repo.Delete(key)
	}

	s.mocks = make([]int32, 0)
}

func (s *Scoped) Done() error {
	if s.IsPending() {
		return ErrScopeNotDone
	}

	return nil
}
