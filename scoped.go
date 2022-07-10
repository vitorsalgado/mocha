package mocha

import (
	"fmt"
	"strings"

	"github.com/vitorsalgado/mocha/core"
)

// Scoped holds references to one or more added mocks allowing users perform operations on them, like enabling/disabling.
type Scoped struct {
	storage core.Storage
	mocks   []*core.Mock
}

func scope(repo core.Storage, mocks []*core.Mock) Scoped {
	return Scoped{storage: repo, mocks: mocks}
}

// Called returns true if all scoped mocks were called at least once.
func (s *Scoped) Called() bool {
	for _, m := range s.mocks {
		if !m.Called() {
			return false
		}
	}

	return true
}

// ListPending returns all mocks that were not called at least once.
func (s *Scoped) ListPending() []core.Mock {
	ret := make([]core.Mock, 0)
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

	s.mocks = make([]*core.Mock, 0)
}

// MustHaveBeenCalled reports a failure if there are still pending mocks.
func (s *Scoped) MustHaveBeenCalled(t core.T) {
	t.Helper()

	if s.IsPending() {
		b := strings.Builder{}
		pending := s.ListPending()
		size := len(pending)

		for _, p := range pending {
			b.WriteString(fmt.Sprintf("	mock: %d %s\n", p.ID, p.Name))
		}

		t.Errorf("\nthere are still %d mocks that were not called.\npending:\n%s", size, b.String())
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
