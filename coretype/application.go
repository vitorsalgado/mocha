package coretype

import (
	"fmt"
	"strings"
	"sync"
)

// TestingT is based on testing.T and is used for assertions.
// See Assert* methods on the application instance.
type TestingT interface {
	Helper()
	Logf(format string, a ...any)
	Errorf(format string, a ...any)
	Cleanup(func())
}

type Assertions interface {
	AssertCalled(t TestingT) bool
	AssertNotCalled(t TestingT) bool
	AssertNumberOfCalls(t TestingT, expected int) bool
}

type MockApp[TMock Mock] interface {
	// Mock(builders ...Builder[TMock]) (*Scope[TMock], error)
	// MustMock(builders ...Builder[TMock]) *Scope[TMock]

	Hits() int
	Enable()
	Disable()
	Clean()
}

type MockBag[TMock Mock, TMockApp MockApp[TMock]] struct {
	rwMutex sync.RWMutex
	app     TMockApp
	store   *MockStore[TMock]
	scopes  []*Scope[TMock]
}

func (a *MockBag[TMock, TMockApp]) Add(builders ...Builder[TMock, TMockApp]) (*Scope[TMock], error) {
	a.rwMutex.Lock()
	defer a.rwMutex.Unlock()

	added := make([]string, len(builders))

	for i, b := range builders {
		mock, err := b.Build(a.app)
		if err != nil {
			return nil, fmt.Errorf("server: error adding mock at index %d.\n%w", i, err)
		}

		mock.Prepare()

		a.store.Save(mock)
		added[i] = mock.GetID()
	}

	scope := NewScope[TMock](a.store, added)
	a.scopes = append(a.scopes, scope)

	return scope, nil
}

func (a *MockBag[TMock, TMockApp]) List() []*Scope[TMock] {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	return a.scopes
}

type BaseApp[TMock Mock, TMockApp MockApp[TMock]] struct {
	mockBag *MockBag[TMock, TMockApp]
	rwMutex sync.RWMutex
}

func NewBaseApp[TMock Mock, TMockApp MockApp[TMock]](app TMockApp, store *MockStore[TMock]) *BaseApp[TMock, TMockApp] {
	return &BaseApp[TMock, TMockApp]{mockBag: &MockBag[TMock, TMockApp]{app: app, store: store}}
}

// Mock adds one or multiple request mocks.
// It returns a Scope instance that allows control of the added mocks and also checks if they were called or not.
//
// Usage:
//
//	scoped := m.MustMock(
//		Get(matcher.URLPath("/test")).
//			Header("test", matcher.StrictEqual("hello")).
//			Query("filter", matcher.StrictEqual("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(txtTemplate, scoped.HasBeenCalled())
func (app *BaseApp[TMock, TMockApp]) Mock(builders ...Builder[TMock, TMockApp]) (*Scope[TMock], error) {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	return app.mockBag.Add(builders...)
}

// MustMock adds one or multiple request mocks.
// It returns a Scope instance that allows control of the added mocks and also checks if they were called or not.
// It fails immediately if any error occurs.
//
// Usage:
//
//	scoped := m.MustMock(
//		Get(matcher.URLPath("/test")).
//			Header("test", matcher.StrictEqual("hello")).
//			Query("filter", matcher.StrictEqual("all")).
//			Reply(reply.Created().PlainText("hello world")))
//
//	assert.True(txtTemplate, scoped.HasBeenCalled())
func (app *BaseApp[TMock, TMockApp]) MustMock(builders ...Builder[TMock, TMockApp]) *Scope[TMock] {
	scoped, err := app.Mock(builders...)
	if err != nil {
		panic(err)
	}

	return scoped
}

// Hits returns the total matched request hits.
func (app *BaseApp[TMock, TMockApp]) Hits() int {
	app.rwMutex.RLock()
	defer app.rwMutex.RUnlock()

	hits := 0

	for _, s := range app.mockBag.List() {
		hits += s.Hits()
	}

	return hits
}

// Enable enables all mocks.
func (app *BaseApp[TMock, TMockApp]) Enable() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, scoped := range app.mockBag.List() {
		scoped.Enable()
	}
}

// Disable disables all mocks.
func (app *BaseApp[TMock, TMockApp]) Disable() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, scoped := range app.mockBag.List() {
		scoped.Disable()
	}
}

// Clean removes all scoped mocks.
func (app *BaseApp[TMock, TMockApp]) Clean() {
	app.rwMutex.Lock()
	defer app.rwMutex.Unlock()

	for _, s := range app.mockBag.List() {
		s.Clean()
	}
}

// --
// Assertions
// --

// AssertCalled asserts that all mocks associated with this instance were called at least once.
func (app *BaseApp[TMock, TMockApp]) AssertCalled(t TestingT) bool {
	t.Helper()

	result := true
	size := 0
	buf := strings.Builder{}

	for i, s := range app.mockBag.List() {
		if s.IsPending() {
			buf.WriteString(fmt.Sprintf("   Scope %d\n", i))
			pending := s.GetPending()
			size += len(pending)

			for _, p := range pending {
				buf.WriteString("    Mock [")
				buf.WriteString(p.GetID())
				buf.WriteString("] ")
				buf.WriteString(p.GetName())
				buf.WriteString("\n")
			}

			result = false
		}
	}

	if !result {
		t.Errorf("\nThere are still %d mocks that were not called.\n  Pending:\n%s",
			size,
			buf.String(),
		)
	}

	return result
}

// AssertNotCalled asserts that all mocks associated with this instance were called at least once.
func (app *BaseApp[TMock, TMockApp]) AssertNotCalled(t TestingT) bool {
	t.Helper()

	result := true
	size := 0
	buf := strings.Builder{}

	for i, s := range app.mockBag.List() {
		if !s.IsPending() {
			buf.WriteString(fmt.Sprintf("   Scope %d\n", i))
			called := s.GetCalled()
			size += len(called)

			for _, p := range called {
				buf.WriteString("    Mock [")
				buf.WriteString(p.GetID())
				buf.WriteString("] ")
				buf.WriteString(p.GetName())
				buf.WriteString("\n")
			}

			result = false
		}
	}

	if !result {
		t.Errorf(
			"\n%d Mock(s) were called at least once when none should be.\n  Called:\n%s",
			size,
			buf.String(),
		)
	}

	return result
}

// AssertNumberOfCalls asserts that the sum of matched request hits
// is equal to the given expected value.
func (app *BaseApp[TMock, TMockApp]) AssertNumberOfCalls(t TestingT, expected int) bool {
	t.Helper()

	hits := app.Hits()

	if hits == expected {
		return true
	}

	t.Errorf("\nExpected %d matched request hits.\n Got %d", expected, hits)

	return false
}
