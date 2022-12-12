package mocha

import (
	"strings"
)

// Scoped holds references to one or more added store allowing users perform operations on them, like enabling/disabling.
type Scoped struct {
	storage mockStore
	mocks   []*Mock
}

func scope(repo mockStore, mocks []*Mock) *Scoped {
	return &Scoped{storage: repo, mocks: mocks}
}

// Get returns a Mock with the given id.
func (s *Scoped) Get(id string) *Mock {
	for _, mock := range s.mocks {
		if mock.ID == id {
			return mock
		}
	}

	return nil
}

// GetAll returns all Mock(s) Scoped in this instance.
func (s *Scoped) GetAll() []*Mock {
	return s.mocks
}

// GetPending returns all Mock(s) that were not called at least once.
func (s *Scoped) GetPending() []*Mock {
	ret := make([]*Mock, 0)
	for _, m := range s.mocks {
		if !m.HasBeenCalled() {
			ret = append(ret, m)
		}
	}

	return ret
}

// GetCalled returns all store that were called.
func (s *Scoped) GetCalled() []*Mock {
	ret := make([]*Mock, 0)
	for _, m := range s.mocks {
		if m.HasBeenCalled() {
			ret = append(ret, m)
		}
	}

	return ret
}

// HasBeenCalled returns true if all Scoped Mock(s) were called at least once.
func (s *Scoped) HasBeenCalled() bool {
	for _, m := range s.mocks {
		if !m.HasBeenCalled() {
			return false
		}
	}

	return true
}

// IsPending returns true when there are one or more store that were not called at least once.
func (s *Scoped) IsPending() bool {
	pending := false
	for _, m := range s.mocks {
		if !m.HasBeenCalled() {
			pending = true
			break
		}
	}

	return pending
}

// Disable scoped store.
// Disabled store will be ignored.
func (s *Scoped) Disable() {
	for _, m := range s.mocks {
		m.Disable()
	}
}

// Enable scoped store.
func (s *Scoped) Enable() {
	for _, m := range s.mocks {
		m.Enable()
	}
}

// Clean all scoped store.
func (s *Scoped) Clean() {
	ids := make([]string, len(s.mocks))

	for i, m := range s.mocks {
		ids[i] = m.ID
	}

	for _, id := range ids {
		s.storage.Delete(id)
	}

	s.mocks = make([]*Mock, 0)
}

// Hits returns the sum of the scoped store calls.
func (s *Scoped) Hits() int {
	total := 0
	for _, m := range s.mocks {
		total += m.Hits()
	}

	return total
}

// AssertCalled reports an error if there are still pending store.
func (s *Scoped) AssertCalled(t TestingT) bool {
	t.Helper()

	if !s.IsPending() {
		return true
	}

	b := strings.Builder{}
	pending := s.GetPending()
	size := len(pending)

	for _, p := range pending {
		b.WriteString("  mock [")
		b.WriteString(p.ID)
		b.WriteString("] ")
		b.WriteString(p.Name)
		b.WriteString("\n")
	}

	t.Errorf("\nthere are still [%d] mocks that were not called.\npending:\n%s", size, b.String())

	return false
}

// AssertNotCalled reports an error if any mock was called.
func (s *Scoped) AssertNotCalled(t TestingT) bool {
	t.Helper()

	if s.IsPending() {
		return true
	}

	b := strings.Builder{}
	called := s.GetCalled()
	size := len(called)

	for _, p := range called {
		b.WriteString("  mock [")
		b.WriteString(p.ID)
		b.WriteString("] ")
		b.WriteString(p.Name)
		b.WriteString("\n")
	}

	t.Errorf("\nthere are [%d] mocks that were called at least once.\ncalled:\n%s", size, b.String())

	return false
}

// AssertCalls asserts that the sum of matched request hits
// is equal to the given expected value.
func (s *Scoped) AssertCalls(t TestingT, expected int) bool {
	t.Helper()

	hits := s.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nexpected [%d] matched request hits. got [%d]", expected, hits)

	return false
}
