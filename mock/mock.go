package mock

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/vitorsalgado/mocha/matcher"
)

func New() *Mock {
	return &Mock{
		ID:          id.Next(),
		Enabled:     true,
		PostActions: make([]PostAction, 0),

		mu: &sync.Mutex{},
	}
}

func (m *Mock) Hit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Hits++
}

func (m *Mock) Called() bool {
	return m.Hits > 0
}

func (m *Mock) Enable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = true
}

func (m *Mock) Disable() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Enabled = false
}

func (m *Mock) Matches(params matcher.Args) (MatchResult, error) {
	weight := 0
	for _, expect := range m.Expectations {
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
		case Expectation[int]:
			matched, w, err = matches(e, params)
		case Expectation[float64]:
			matched, w, err = matches(e, params)
		case Expectation[bool]:
			matched, w, err = matches(e, params)
		case Expectation[map[string]any]:
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
