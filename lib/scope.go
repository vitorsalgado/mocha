package lib

import (
	"reflect"
	"strings"
)

// Scope holds keeps a reference to a group of mocks.
// With a Scope instance, it is possible to verify if one or a group of mocks were called,
// how many times they were called and so on.
type Scope[TMock Mock] struct {
	store *MockStore[TMock]
	ids   map[string]struct{}
}

func NewScope[TMock Mock](store *MockStore[TMock], ids []string) *Scope[TMock] {
	datum := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		datum[id] = struct{}{}
	}

	return &Scope[TMock]{store, datum}
}

// Get returns a Mock with the given ID.
func (s *Scope[TMock]) Get(id string) TMock {
	_, ok := s.ids[id]
	if !ok {
		var result TMock
		return result
	}

	return s.store.Get(id)
}

// GetAll returns all Mock instances kept in this Scope.
func (s *Scope[TMock]) GetAll() []TMock {
	mocks := make([]TMock, 0, len(s.ids))

	for id := range s.ids {
		if m := s.store.Get(id); !reflect.ValueOf(m).IsZero() {
			mocks = append(mocks, m)
		}
	}

	return mocks
}

// GetPending returns all Mock instances that were not called at least once.
func (s *Scope[TMock]) GetPending() []TMock {
	mocks := make([]TMock, 0, len(s.ids))

	for _, m := range s.GetAll() {
		if !m.HasBeenCalled() {
			mocks = append(mocks, m)
		}
	}

	return mocks
}

// GetCalled returns all Mock instances that were called.
func (s *Scope[TMock]) GetCalled() []TMock {
	mocks := make([]TMock, 0, len(s.ids))

	for _, m := range s.GetAll() {
		if m.HasBeenCalled() {
			mocks = append(mocks, m)
		}
	}

	return mocks
}

// HasBeenCalled returns true if all Scope Mock instances were called at least once.
func (s *Scope[TMock]) HasBeenCalled() bool {
	for _, m := range s.GetAll() {
		if !m.HasBeenCalled() {
			return false
		}
	}

	return true
}

// IsPending returns true when there are one or more Mock instances that were not called at least once.
func (s *Scope[TMock]) IsPending() bool {
	for _, m := range s.GetAll() {
		if !m.HasBeenCalled() {
			return true
		}
	}

	return false
}

// Disable scoped store.
// Disabled Mock will be ignored.
func (s *Scope[TMock]) Disable() {
	for _, m := range s.GetAll() {
		m.Disable()
	}
}

// Enable all Mock instances kept in this Scope.
func (s *Scope[TMock]) Enable() {
	for _, m := range s.GetAll() {
		m.Enable()
	}
}

// Delete removes a by ID, as long as it is scoped by this instance.
func (s *Scope[TMock]) Delete(id string) bool {
	_, ok := s.ids[id]
	if !ok {
		return false
	}

	s.store.Delete(id)
	delete(s.ids, id)

	return true
}

// Clean all scoped Mock instances.
func (s *Scope[TMock]) Clean() {
	for id := range s.ids {
		s.store.Delete(id)
		delete(s.ids, id)
	}
}

// Hits returns the sum of the Scope store calls.
func (s *Scope[TMock]) Hits() int {
	total := 0
	for _, m := range s.GetAll() {
		total += m.Hits()
	}

	return total
}

// AssertCalled reports an error if there are still pending Mock instances.
func (s *Scope[TMock]) AssertCalled(t TestingT) bool {
	t.Helper()

	if !s.IsPending() {
		return true
	}

	b := strings.Builder{}
	pending := s.GetPending()
	size := len(pending)

	for _, p := range pending {
		b.WriteString("   Mock [")
		b.WriteString(p.GetID())
		b.WriteString("] ")
		b.WriteString(p.GetName())
		b.WriteString("\n")
	}

	t.Errorf("\nThere are still %d mocks that were not called.\nPending:\n%s", size, b.String())

	return false
}

// AssertNotCalled reports an error if any mock was called.
func (s *Scope[TMock]) AssertNotCalled(t TestingT) bool {
	t.Helper()

	if s.IsPending() {
		return true
	}

	b := strings.Builder{}
	called := s.GetCalled()
	size := len(called)

	for _, p := range called {
		b.WriteString("  Mock [")
		b.WriteString(p.GetID())
		b.WriteString("] ")
		b.WriteString(p.GetName())
		b.WriteString("\n")
	}

	t.Errorf("\n%d Mocks were called at least once when none should be.\nCalled:\n%s", size, b.String())

	return false
}

// AssertNumberOfCalls asserts that the sum of matched request hits
// is equal to the given expected value.
func (s *Scope[TMock]) AssertNumberOfCalls(t TestingT, expected int) bool {
	t.Helper()

	hits := s.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nExpected %d matched request hits.\nGot %d", expected, hits)

	return false
}
