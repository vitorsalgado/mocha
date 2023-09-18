package dzstd

import (
	"context"
	"strings"
	"sync"
)

// Scope keeps a reference to a group of mocks.
// With a Scope instance, it is possible to verify if one or a group of mocks were called,
// how many times they were called and so on.
type Scope[TMock Mock] struct {
	mocks      map[string]TMock
	del        map[string]struct{}
	repository MockRepository[TMock]
	rwmu       sync.RWMutex
}

func NewScope[TMock Mock](repository MockRepository[TMock], mocks []TMock) *Scope[TMock] {
	scoped := &Scope[TMock]{mocks: make(map[string]TMock, len(mocks)), del: make(map[string]struct{}), repository: repository}
	for _, m := range mocks {
		scoped.mocks[m.GetID()] = m
	}

	return scoped
}

func (s *Scope[TMock]) IDs() []string {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	ids := make([]string, 0, len(s.mocks))
	for _, id := range s.mocks {
		ids = append(ids, id.GetID())
	}

	return ids
}

func (s *Scope[TMock]) Exists(id string) bool {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	_, exists := s.mocks[id]
	return exists
}

// Get returns a Mock with the given ID.
func (s *Scope[TMock]) Get(id string) TMock {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	return s.mocks[id]
}

// GetAll returns all Mock instances kept in this Scope.
func (s *Scope[TMock]) GetAll() []TMock {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	mocks := make([]TMock, len(s.mocks))
	i := 0

	for _, m := range s.mocks {
		mocks[i] = m
		i++
	}

	return mocks
}

// GetPending returns all Mock instances that were not called at least once.
func (s *Scope[TMock]) GetPending() []TMock {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	mocks := make([]TMock, 0, len(s.mocks))

	for _, m := range s.mocks {
		if !m.HasBeenCalled() {
			mocks = append(mocks, m)
		}
	}

	return mocks
}

// GetCalled returns all Mock instances that were called.
func (s *Scope[TMock]) GetCalled() []TMock {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	mocks := make([]TMock, 0, len(s.mocks))

	for _, m := range s.mocks {
		if m.HasBeenCalled() {
			mocks = append(mocks, m)
		}
	}

	return mocks
}

// HasBeenCalled returns true if all Scope Mock instances were called at least once.
func (s *Scope[TMock]) HasBeenCalled() bool {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

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
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	_, ok := s.mocks[id]
	if !ok {
		return false
	}

	delete(s.mocks, id)
	s.del[id] = struct{}{}

	return true
}

// Clean all scoped Mock instances.
func (s *Scope[TMock]) Clean(ctx context.Context) error {
	for id := range s.mocks {
		s.Delete(id)
	}

	return s.Sync(ctx)
}

// Hits returns the sum of the Scope store calls.
func (s *Scope[TMock]) Hits() int64 {
	total := int64(0)
	for _, m := range s.GetAll() {
		total += m.Hits()
	}

	return total
}

func (s *Scope[TMock]) Sync(ctx context.Context) error {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	mocks := make([]TMock, 0, len(s.mocks))
	for _, v := range s.mocks {
		mocks = append(mocks, v)
	}

	if err := s.repository.Save(ctx, mocks...); err != nil {
		return err
	}

	ids := make([]string, 0, len(s.mocks))
	for k := range s.del {
		ids = append(ids, k)
	}

	for _, id := range ids {
		err := s.repository.Delete(ctx, id)
		if err != nil {
			return err
		}

		delete(s.del, id)
	}

	return nil
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
func (s *Scope[TMock]) AssertNumberOfCalls(t TestingT, expected int64) bool {
	t.Helper()

	hits := s.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nExpected %d matched request hits.\nGot %d", expected, hits)

	return false
}
