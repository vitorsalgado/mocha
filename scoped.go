package mocha

import (
	"errors"

	"github.com/vitorsalgado/mocha/mock"
)

// Scoped holds references to one or more added mocks allowing users perform operations on them, like enabling/disabling.
type Scoped struct {
	storage mock.Storage
	mocks   []*mock.Mock
}

var (
	// ErrScopeNotDone is returned when scope was not called.
	ErrScopeNotDone = errors.New("there are still mocks that were not called")
)

func scope(repo mock.Storage, mocks []*mock.Mock) Scoped {
	return Scoped{storage: repo, mocks: mocks}
}

// IsDone returns true if all scoped mocks were called at least once.
func (s *Scoped) IsDone() bool {
	for _, m := range s.mocks {
		if !m.Called() {
			return false
		}
	}

	return true
}

// Pending returns all mocks that were not called at least once.
func (s *Scoped) Pending() []mock.Mock {
	ret := make([]mock.Mock, 0)
	for _, m := range s.mocks {
		if !m.Called() {
			ret = append(ret, *m)
		}
	}

	return ret
}

// IsPending returns true when there are one or more mocks that were not called at least once.
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

// Disable scoped mocks.
// Disabled mocks will be ignored.
func (s *Scoped) Disable() {
	for _, m := range s.mocks {
		m.Disable()
	}
}

// Enable scoped mocks.
func (s *Scoped) Enable() {
	for _, m := range s.mocks {
		m.Enable()
	}
}

// Clean all scoped mocks.
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

// MustBeDone panic if there are still pending mocks.
func (s *Scoped) MustBeDone() {
	if s.IsPending() {
		panic(ErrScopeNotDone)
	}
}

// Hits returns the sum of the scoped mocks calls.
func (s *Scoped) Hits() int {
	total := 0
	for _, m := range s.mocks {
		total += m.Hits
	}

	return total
}
