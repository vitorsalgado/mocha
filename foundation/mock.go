package foundation

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/vitorsalgado/mocha/v3/matcher"
)

var (
	_ Mock = (*BaseMock)(nil)
)

type Mock interface {
	GetID() string
	GetName() string
	GetPriority() int
	GetSource() string
	Inc()
	Dec()
	Hits() int
	HasBeenCalled() bool
	IsEnabled() bool
	Enable()
	Disable()
	Prepare()
}

// Builder describes a Mock builder.
type Builder[TMock Mock, TMockApp MockApp[TMock]] interface {
	Build(app TMockApp) (TMock, error)
}

type RequestMatcher[TValueIn any] interface {
	GetExpectations() []*Expectation[TValueIn]
}

// BaseMock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is the core entity of this project, and most features work based on it.
type BaseMock struct {
	// ID is the unique identifier of a Mock
	ID string

	// Name describes the mock. Useful for debugging.
	Name string

	// Priority sets the priority of a Mock.
	Priority int

	// Enabled indicates if the Mock is enabled or disabled.
	// Only enabled mocks are considered during the request matching phase.
	Enabled bool

	// Source describes the source of the mock. E.g.: if it was built from a file,
	// it will contain the filename.
	Source string

	// Delay sets the duration to delay serving the mocked response.
	Delay time.Duration

	after []matcher.OnAfterMockServed
	mu    sync.RWMutex
	hits  int
}

// NewMock returns a new Mock with default values set.
func NewMock() *BaseMock {
	return &BaseMock{
		ID:      uuid.New().String(),
		Enabled: true,
	}
}

func (m *BaseMock) GetID() string {
	return m.ID
}

func (m *BaseMock) GetName() string {
	return m.Name
}

func (m *BaseMock) GetPriority() int {
	return m.Priority
}

func (m *BaseMock) GetSource() string {
	return m.Source
}

func (m *BaseMock) IsEnabled() bool {
	return m.Enabled
}

// Inc increment one Mock call.
func (m *BaseMock) Inc() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits++
}

// Dec reduce one Mock call.
func (m *BaseMock) Dec() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hits--
}

// Hits returns the amount of time this Mock was matched to a request and served.
func (m *BaseMock) Hits() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits
}

// HasBeenCalled checks if the Mock was called at least once.
func (m *BaseMock) HasBeenCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hits > 0
}

// Enable enables the Mock.
// The Mock will be eligible to be matched.
func (m *BaseMock) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = true
}

// Disable disables the Mock.
// The Mock will not be eligible to be matched.
func (m *BaseMock) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = false
}

func (m *BaseMock) Prepare() {
}

// Build allows users to use the Mock itself as a HTTPMockBuilder.
func (m *BaseMock) Build() (*BaseMock, error) {
	return m, nil
}

type MockStore[TMock Mock] struct {
	data    []TMock
	rwMutex sync.RWMutex
}

// NewStore returns Mock MockStore implementation.
func NewStore[TMock Mock]() *MockStore[TMock] {
	return &MockStore[TMock]{data: make([]TMock, 0)}
}

func (s *MockStore[TMock]) Save(mock TMock) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = append(s.data, mock)

	sort.SliceStable(s.data, func(a, b int) bool {
		return s.data[a].GetPriority() < s.data[b].GetPriority()
	})
}

func (s *MockStore[TMock]) Get(id string) TMock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	for _, datum := range s.data {
		if datum.GetID() == id {
			return datum
		}
	}

	var result TMock
	return result
}

func (s *MockStore[TMock]) GetEligible() []TMock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	mocks := make([]TMock, 0, len(s.data))

	for _, mock := range s.data {
		if mock.IsEnabled() {
			mocks = append(mocks, mock)
		}
	}

	return mocks
}

func (s *MockStore[TMock]) GetAll() []TMock {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	return s.data
}

func (s *MockStore[TMock]) Delete(id string) {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	index := -1
	for i, m := range s.data {
		if m.GetID() == id {
			index = i
			break
		}
	}

	s.data = s.data[:index+copy(s.data[index:], s.data[index+1:])]
}

func (s *MockStore[TMock]) DeleteExternal() {
	ids := make([]string, 0, len(s.data))

	for _, m := range s.data {
		if len(m.GetSource()) > 0 {
			ids = append(ids, m.GetID())
		}
	}

	for _, id := range ids {
		s.Delete(id)
	}
}

func (s *MockStore[TMock]) DeleteAll() {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	s.data = nil
	s.data = make([]TMock, 0)
}

// Expectation holds metadata related to one http.Request Matcher.
type Expectation[TValueIn any] struct {
	// Target is an optional metadata that describes the target of the matcher.
	// Example: the target could have the "header", meaning that the matcher will be applied to one request misc.Header
	Target Target

	Key string

	// Matcher associated with this Expectation.
	Matcher matcher.Matcher

	// ValueSelector will extract the http.Request or a specific field of it and feed it to the associated Matcher.
	ValueSelector func(TValueIn) any

	// Weight of this Expectation.
	Weight Weight
}

// Weight helps to detect the closest mock match.
type Weight int8

// Enum of Weight.
const (
	WeightNone Weight = iota
	WeightLow  Weight = iota * 2
	WeightVeryLow
	WeightRegular
	WeightHigh
)

// MatchResult holds information related to a matching operation.
type MatchResult struct {
	// Details is the list of non-matches messages.
	Details []MismatchDetail

	// Weight for the matcher. It helps determine the closest match.
	Weight int

	// Pass indicates whether it matched or not.
	Pass bool
}

// MismatchDetail gives more context about why a matcher did not match.
type MismatchDetail struct {
	MatchersName string
	Target       Target
	Key          string
	Result       *matcher.Result
	Err          error
}

func matchExpectation[VS any](e *Expectation[VS], value any) (result *matcher.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: matcher=%s. %v", e.Matcher.Name(), r)
			return
		}
	}()

	result, err = e.Matcher.Match(value)
	if err != nil {
		err = fmt.Errorf("%s: error while matching. %w", e.Matcher.Name(), err)
	}

	return
}

type Target int8

// Target constants to help debug unmatched requests.
const (
	TargetRequest Target = iota
	TargetScheme
	TargetMethod
	TargetURL
	TargetHeader
	TargetQuery
	TargetBody
	TargetForm
)

func (mt Target) String() string {
	switch mt {
	case TargetRequest:
		return "request"
	case TargetScheme:
		return "scheme"
	case TargetMethod:
		return "method"
	case TargetURL:
		return "url"
	case TargetHeader:
		return "header"
	case TargetQuery:
		return "query"
	case TargetBody:
		return "body"
	case TargetForm:
		return "form"
	default:
		return ""
	}
}
