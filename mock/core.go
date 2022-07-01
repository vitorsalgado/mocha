package mock

import (
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/matcher"
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
		Params   params.Params
	}

	PostAction interface {
		Run(args PostActionArgs) error
	}

	Response struct {
		Status  int
		Header  http.Header
		Cookies []*http.Cookie
		Body    io.Reader
		Delay   time.Duration
		Err     error
		Mappers []ResponseMapper
	}

	Reply interface {
		Build(*http.Request, *Mock, params.Params) (*Response, error)
	}

	ResponseMapperArgs struct {
		Request    *http.Request
		Parameters params.Params
	}

	ResponseMapper func(res *Response, args ResponseMapperArgs) error

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
)

type autoID struct {
	id int32
}

func (i *autoID) Next() int {
	atomic.AddInt32(&i.id, 1)
	return int(i.id)
}

var id = autoID{}
