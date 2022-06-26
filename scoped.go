package mocha

import (
	"errors"

	"github.com/vitorsalgado/mocha/mock"
)

type Scoped struct {
	storage mock.Storage
	mocks   []*mock.Mock
}

var (
	ErrScopeNotDone = errors.New("there are still mocks that were not called")
)

func Scope(repo mock.Storage, mocks []*mock.Mock) *Scoped {
	return &Scoped{storage: repo, mocks: mocks}
}

func (s *Scoped) IsDone() bool {
	for _, m := range s.mocks {
		if !m.Called() {
			return false
		}
	}

	return true
}

func (s *Scoped) Pending() []mock.Mock {
	ret := make([]mock.Mock, 0)
	for _, m := range s.mocks {
		if !m.Called() {
			ret = append(ret, *m)
		}
	}

	return ret
}

func (s *Scoped) IsPending() bool {
	pending := false
	for _, m := range s.mocks {
		if !m.Called() {
			pending = true
			break
		}
	}

	return pending
}

func (s *Scoped) Disable() {
	for _, m := range s.mocks {
		m.Disable()
	}
}

func (s *Scoped) Enable() {
	for _, m := range s.mocks {
		m.Enable()
	}
}

func (s *Scoped) Clean() {
	ids := make([]int, len(s.mocks))

	for i, m := range s.mocks {
		ids[i] = m.ID
	}

	for _, id := range ids {
		s.storage.Delete(id)
	}

	s.mocks = make([]*mock.Mock, 0)
}

func (s *Scoped) Done() error {
	if s.IsPending() {
		return ErrScopeNotDone
	}

	return nil
}
