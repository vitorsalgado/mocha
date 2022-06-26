package mock

import (
	"fmt"
	"github.com/vitorsalgado/mocha/internal"
	"github.com/vitorsalgado/mocha/matcher"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"time"
)

type (
	Mock struct {
		ID           int32
		Name         string
		Priority     int
		Expectations []any
		Reply        Reply
		Hits         int
	}

	Storage interface {
		Save(mock *Mock)
		FetchSorted() []Mock
		FetchByID(id int32) Mock
		FetchByIDs(ids []int32) []Mock
		Pending(ids []int32) []Mock
		Delete(id int32)
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
		Body    []byte
		Delay   time.Duration
		Err     error
	}

	Responder func(io.Writer, *http.Request, *Mock) error

	Reply interface {
		Err() error
		Build(*http.Request, *Mock) (*Response, error)
	}
)

func New() *Mock {
	return &Mock{}
}

func (m *Mock) Hit() {
	m.Hits++
}

func (m *Mock) Called() bool {
	return m.Hits > 0
}

func (m *Mock) Matches(ctx matcher.Params) (MatchResult, error) {
	weight := 0
	for _, expect := range m.Expectations {
		switch e := expect.(type) {
		case Expectation[string]:
			if res, err := e.Matcher(e.ValuePicker(ctx.RequestInfo), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[url.URL]:
			if res, err := e.Matcher(e.ValuePicker(ctx.RequestInfo), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[*http.Request]:
			if res, err := e.Matcher(e.ValuePicker(ctx.RequestInfo), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[any]:
			if res, err := e.Matcher(e.ValuePicker(ctx.RequestInfo), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		default:
			return MatchResult{IsMatch: false, Weight: weight}, fmt.Errorf("unhandled matcher type %s", reflect.TypeOf(e))
		}
	}

	return MatchResult{IsMatch: true}, nil
}

type InMemoMockStore struct {
	data map[int32]Mock
	ids  internal.ID
}

func NewMockStorage() Storage {
	return &InMemoMockStore{data: make(map[int32]Mock)}
}

func (repo *InMemoMockStore) Save(mock *Mock) {
	if mock.ID == 0 {
		mock.ID = repo.ids.Next()
	}

	repo.data[mock.ID] = *mock
}

func (repo *InMemoMockStore) FetchByID(id int32) Mock {
	return repo.data[id]
}

func (repo *InMemoMockStore) FetchSorted() []Mock {
	size := len(repo.data)
	mocks := make([]Mock, size)
	i := 0

	for _, mock := range repo.data {
		mocks[i] = mock
		i++
	}

	sort.SliceStable(mocks, func(a, b int) bool {
		return mocks[a].Priority < mocks[b].Priority
	})

	return mocks
}

func (repo *InMemoMockStore) FetchByIDs(ids []int32) []Mock {
	size := len(ids)
	arr := make([]Mock, 0, size)

	for _, id := range ids {
		arr = append(arr, repo.FetchByID(id))
	}

	return arr
}

func (repo *InMemoMockStore) Pending(ids []int32) []Mock {
	r := make([]Mock, 0)

	for _, mock := range repo.FetchByIDs(ids) {
		if !mock.Called() {
			r = append(r, mock)
		}
	}

	return r
}

func (repo *InMemoMockStore) Delete(id int32) {
	delete(repo.data, id)
}

func (repo *InMemoMockStore) Flush() {
	for key := range repo.data {
		delete(repo.data, key)
	}
}
