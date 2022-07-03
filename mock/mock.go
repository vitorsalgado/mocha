package mock

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/vitorsalgado/mocha/matcher"
)

// New returns a new Mock with default values set.
func New() *Mock {
	return &Mock{
		ID:                id.next(),
		Enabled:           true,
		Expectations:      make([]any, 0),
		AfterExpectations: make([]any, 0),
		PostActions:       make([]PostAction, 0),

		mu: &sync.Mutex{},
	}
}

// Hit notify that the Mock was called.
func (m *Mock) Hit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits++
}

// Dec reduce one Mock call.
func (m *Mock) Dec() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits--
}

// Called checks if the Mock was called at least once.
func (m *Mock) Called() bool {
	return m.Hits > 0
}

// Enable enables the Mock.
// The Mock will be eligible to be matched.
func (m *Mock) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = true
}

// Disable disables the Mock.
// The Mock will not be eligible to be matched.
func (m *Mock) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = false
}

// Matches checks if current Mock matches against a list of expectations.
func (m *Mock) Matches(params matcher.Args, expectations []any) (MatchResult, error) {
	weight := 0
	for _, expect := range expectations {
		var matched bool
		var err error
		var w int

		switch e := expect.(type) {
		default:
			return MatchResult{IsMatch: false, Weight: weight},
				fmt.Errorf("unhandled matcher type %s", reflect.TypeOf(e))

		case Expectation[any]:
			matched, w, err = matches(e, params)
		case Expectation[string]:
			matched, w, err = matches(e, params)
		case Expectation[float64]:
			matched, w, err = matches(e, params)
		case Expectation[bool]:
			matched, w, err = matches(e, params)
		case Expectation[map[string]any]:
			matched, w, err = matches(e, params)
		case Expectation[map[string]string]:
			matched, w, err = matches(e, params)
		case Expectation[map[string][]string]:
			matched, w, err = matches(e, params)
		case Expectation[[]any]:
			matched, w, err = matches(e, params)
		case Expectation[url.URL]:
			matched, w, err = matches(e, params)
		case Expectation[*http.Request]:
			matched, w, err = matches(e, params)
		case Expectation[url.Values]:
			matched, w, err = matches(e, params)
		}

		if err != nil || !matched {
			return MatchResult{IsMatch: false, Weight: weight}, err
		}

		weight += w
	}

	return MatchResult{IsMatch: true}, nil
}

func matches[V any](e Expectation[V], params matcher.Args) (bool, int, error) {
	res, err := e.Matcher(e.ValuePicker(params.RequestInfo), params)
	return res, e.Weight, err
}
