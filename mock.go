package mocha

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"

	"github.com/vitorsalgado/mocha/internal"
)

type (
	Mock struct {
		ID           int32
		Name         string
		Priority     int
		Expectations []any
		ResFn        ResponseDelegate
		Hits         int
	}

	MockRepository interface {
		Save(mock *Mock)
		FetchSorted() []Mock
		GetByID(id int32) Mock
		GetByIDs(ids []int32) []Mock
		Pending(ids []int32) []Mock
		Delete(id int32)
		Flush()
	}

	RequestPicker[V any] func(r *MockRequest) V

	Expectation[V any] struct {
		Matcher Matcher[V]
		Pick    RequestPicker[V]
		Weight  int
		Name    string
	}

	MatchResult struct {
		NonMatched []string
		Weight     int
		IsMatch    bool
	}
)

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Hit() {
	m.Hits++
}

func (m *Mock) Called() bool {
	return m.Hits > 0
}

func (m *Mock) Matches(ctx MatcherParams) (MatchResult, error) {
	weight := 0

	for _, expect := range m.Expectations {
		switch e := expect.(type) {
		case Expectation[string]:
			if res, err := e.Matcher(e.Pick(ctx.Req), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[url.URL]:
			if res, err := e.Matcher(e.Pick(ctx.Req), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[*http.Request]:
			if res, err := e.Matcher(e.Pick(ctx.Req), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		case Expectation[any]:
			if res, err := e.Matcher(e.Pick(ctx.Req), ctx); err != nil || !res {
				return MatchResult{IsMatch: false, Weight: weight}, err
			}
			weight = weight + e.Weight
		default:
			return MatchResult{IsMatch: false, Weight: weight}, fmt.Errorf("unhandled matcher type %s", reflect.TypeOf(e))
		}
	}

	return MatchResult{IsMatch: true}, nil
}

type InMemoryMockRepository struct {
	data map[int32]Mock
	ids  internal.ID
}

func NewMockRepository() MockRepository {
	return &InMemoryMockRepository{data: make(map[int32]Mock)}
}

func (repo *InMemoryMockRepository) Save(mock *Mock) {
	if mock.ID == 0 {
		mock.ID = repo.ids.Next()
	}

	repo.data[mock.ID] = *mock
}

func (repo *InMemoryMockRepository) GetByID(id int32) Mock {
	return repo.data[id]
}

func (repo *InMemoryMockRepository) FetchSorted() []Mock {
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

func (repo *InMemoryMockRepository) GetByIDs(ids []int32) []Mock {
	size := len(ids)
	arr := make([]Mock, 0, size)

	for _, id := range ids {
		arr = append(arr, repo.GetByID(id))
	}

	return arr
}

func (repo *InMemoryMockRepository) Pending(ids []int32) []Mock {
	r := make([]Mock, 0)

	for _, mock := range repo.GetByIDs(ids) {
		if !mock.Called() {
			r = append(r, mock)
		}
	}

	return r
}

func (repo *InMemoryMockRepository) Delete(id int32) {
	delete(repo.data, id)
}

func (repo *InMemoryMockRepository) Flush() {
	for key := range repo.data {
		delete(repo.data, key)
	}
}
