package dzstd

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"

	"github.com/vitorsalgado/mocha/v3/internal/encutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

var (
	_ Mock = (*BaseMock)(nil)
)

type Mock interface {
	json.Marshaler

	GetID() string
	GetName() string
	GetPriority() int
	GetSource() string
	Inc()
	Dec()
	Hits() int64
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

// BaseMock holds metadata and expectations to be matched against HTTP requests in order to serve mocked responses.
// This is the core entity of this project, and most features work based on it.
type BaseMock struct {
	hits    atomic.Int64
	enabled bool
	rwmu    sync.RWMutex

	// ID is the unique identifier of a Mock
	ID string

	// Name describes the mock. Useful for debugging.
	Name string

	// Priority sets the priority of a Mock.
	Priority int

	// Source describes the source of the mock. E.g.: if it was built from a file,
	// it will contain the filename.
	Source string

	// Delay sets the duration to delay serving the mocked response.
	Delay Delay
}

// NewMock returns a new Mock with default values set.
func NewMock() *BaseMock {
	mock := &BaseMock{
		ID:      uuid.New().String(),
		enabled: true,
	}

	return mock
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
	m.rwmu.RLock()
	defer m.rwmu.RUnlock()
	return m.enabled
}

// Inc increment one Mock call.
func (m *BaseMock) Inc() {
	m.hits.Add(1)
}

// Dec reduce one Mock call.
func (m *BaseMock) Dec() {
	m.hits.Add(-1)
}

// Hits returns the amount of time this Mock was matched to a request and served.
func (m *BaseMock) Hits() int64 {
	return m.hits.Load()
}

// HasBeenCalled checks if the Mock was called at least once.
func (m *BaseMock) HasBeenCalled() bool {
	return m.hits.Load() > 0
}

// Enable enables the Mock.
// The Mock will be eligible to be matched.
func (m *BaseMock) Enable() {
	m.rwmu.Lock()
	defer m.rwmu.Unlock()
	m.enabled = true
}

// Disable disables the Mock.
// The Mock will not be eligible to be matched.
func (m *BaseMock) Disable() {
	m.rwmu.Lock()
	defer m.rwmu.Unlock()
	m.enabled = false
}

func (m *BaseMock) Switch(state bool) {
	m.rwmu.Lock()
	defer m.rwmu.Unlock()
	m.enabled = state
}

func (m *BaseMock) Prepare() {
}

// Build allows users to use the Mock itself as a HTTPMockBuilder.
func (m *BaseMock) Build() (*BaseMock, error) {
	return m, nil
}

func (m *BaseMock) MarshalJSON() ([]byte, error) {
	return nil, nil
}

// Expectation holds metadata related to one http.Request Matcher.
type Expectation[TValueIn any] struct {
	Target int

	Key string

	// TargetDescription is an optional metadata that describes the target of the matcher.
	// Eg.: Header(Content-Type)
	TargetDescription string

	// Matcher associated with this Expectation.
	Matcher matcher.Matcher

	// ValueSelector will extract the http.Request or a specific field of it and feed it to the associated Matcher.
	ValueSelector func(context.Context, TValueIn) any

	// Weight of this Expectation.
	Weight Weight
}

// Weight helps to detect the closest mock match.
type Weight int8

// Enum of Weight.
const (
	WeightNone    Weight = iota
	WeightVeryLow Weight = iota * 2
	WeightLow
	WeightRegular
	WeightHigh
)

type Results struct {
	Buf []string
}

func (d *Results) Append(v string) {
	d.Buf = append(d.Buf, v)
}

func (d *Results) AppendList(sep string, v ...string) {
	d.Buf = append(d.Buf, encutil.Join(sep, v...))
}

func (d *Results) Len() int {
	return len(d.Buf)
}

func (d *Results) String() string {
	return encutil.Join("\n", d.Buf...)
}
