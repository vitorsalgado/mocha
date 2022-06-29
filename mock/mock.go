package mock

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"time"

	"github.com/vitorsalgado/mocha/matcher"
	"github.com/vitorsalgado/mocha/params"
)

type (
	Mock struct {
		ID           int
		Name         string
		Priority     int
		Expectations []any
		Reply        Reply
		Hits         int
		Enabled      bool
		PostActions  []PostAction

		mu *sync.Mutex
	}

	PostActionArgs struct {
		Request  *http.Request
		Response *Response
		Mock     *Mock
		Params   *params.Params
	}

	PostAction interface {
		Run(args PostActionArgs) error
	}

	Storage interface {
		Save(mock *Mock)
		FetchEligible() []*Mock
		FetchAll() []*Mock
		Delete(id int)
		Flush()
	}

	ExpectationValuePicker[V any] func(r *matcher.RequestInfo) V

	Expectation[V any] struct {
		Name        string
		Matcher     matcher.Matcher[V]
		ValuePicker ExpectationValuePicker[V]
		Weight      int
	}

	MatchResult struct {
		NonMatched []string
		Weight     int
		IsMatch    bool
	}

	Response struct {
		Status  int
		Header  http.Header
		Cookies []*http.Cookie
		Body    io.Reader
		Delay   time.Duration
		Err     error
	}

	Reply interface {
		Err() error
		Build(*http.Request, *Mock, *params.Params) (*Response, error)
	}
)

var id = autoID{}

func New() *Mock {
	return &Mock{ID: id.Next(), Enabled: true, PostActions: make([]PostAction, 0), mu: &sync.Mutex{}}
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

func (m *Mock) Matches(params matcher.Params) (MatchResult, error) {
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

func matches[V any](e Expectation[V], params matcher.Params) (bool, int, error) {
	res, err := e.Matcher(e.ValuePicker(params.RequestInfo), params)
	return res, e.Weight, err
}
