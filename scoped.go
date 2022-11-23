package mocha

import (
	"fmt"
	"strings"
)

// Scoped holds references to one or more added mocks allowing users perform operations on them, like enabling/disabling.
type Scoped struct {
	storage storage
	mocks   []*Mock
}

func scope(repo storage, mocks []*Mock) *Scoped {
	return &Scoped{storage: repo, mocks: mocks}
}

// Get returns a mock with the provided id.
func (s *Scoped) Get(id int) *Mock {
	for _, mock := range s.mocks {
		if mock.ID == id {
			return mock
		}
	}

	return nil
}

// ListAll returns all mocks scoped in this instance.
func (s *Scoped) ListAll() []*Mock {
	return s.mocks
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
func (s *Scoped) ListPending() []*Mock {
	ret := make([]*Mock, 0)
	for _, m := range s.mocks {
		if !m.Called() {
			ret = append(ret, m)
		}
	}

	return ret
}

// ListCalled returns all mocks that were called.
func (s *Scoped) ListCalled() []*Mock {
	ret := make([]*Mock, 0)
	for _, m := range s.mocks {
		if m.Called() {
			ret = append(ret, m)
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

	s.mocks = make([]*Mock, 0)
}

// AssertCalled reports an error if there are still pending mocks.
func (s *Scoped) AssertCalled(t TestingT) bool {
	t.Helper()

	if s.IsPending() {
		b := strings.Builder{}
		pending := s.ListPending()
		size := len(pending)

		for _, p := range pending {
			b.WriteString(fmt.Sprintf("	mock: %d %s\n", p.ID, p.Name))
		}

		t.Errorf("\nthere are still %d mocks that were not called.\npending:\n%s", size, b.String())

		return false
	}

	return true
}

// AssertNotCalled reports an error if any mock was called.
func (s *Scoped) AssertNotCalled(t TestingT) bool {
	t.Helper()

	if !s.IsPending() {
		b := strings.Builder{}
		called := s.ListCalled()
		size := len(called)

		for _, p := range called {
			b.WriteString(fmt.Sprintf("	mock: %d %s\n", p.ID, p.Name))
		}

		t.Errorf("\nthere are %d mocks that were called at least once.\ncalled:\n%s", size, b.String())

		return false
	}

	return true
}

// Hits returns the sum of the scoped mocks calls.
func (s *Scoped) Hits() int {
	total := 0
	for _, m := range s.mocks {
		total += m.Hits()
	}

	return total
}
