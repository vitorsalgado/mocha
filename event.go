package mocha

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/worker"
	"github.com/vitorsalgado/mocha/v3/mod"
)

// Event Types.
var (
	EventOnRequest           = reflect.TypeOf(&OnRequest{})
	EventOnRequestMatched    = reflect.TypeOf(&OnRequestMatch{})
	EventOnRequestNotMatched = reflect.TypeOf(&OnRequestNotMatched{})
	EventOnError             = reflect.TypeOf(&OnError{})
)

// OnRequest event is triggered every time a request arrives at the mock handler.
type OnRequest struct {
	Request   *mod.EvtReq
	StartedAt time.Time
}

// OnRequestMatch event is triggered when a mock is found for a request.
type OnRequestMatch struct {
	Request            *mod.EvtReq
	ResponseDefinition mod.EvtRes
	Mock               mod.EvtMk
	Elapsed            time.Duration
	Body               any
}

// OnRequestNotMatched event is triggered when no mocks are found for a request.
type OnRequestNotMatched struct {
	Request *mod.EvtReq
	Result  mod.EvtResult
}

// OnError event is triggered when an error occurs during request matching.
type OnError struct {
	Request *mod.EvtReq
	Err     error
}

type eventListener struct {
	w     *worker.Worker
	queue worker.Queue
	jobs  map[worker.JobType][]func(e any)
	mu    sync.Mutex
}

func newEvents() *eventListener {
	h := &eventListener{}
	h.jobs = map[worker.JobType][]func(e any){}
	h.w = &worker.Worker{Jobs: h.jobs}

	return h
}

// StartListening starts background event listener.
func (h *eventListener) StartListening(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.w.Started.Load() {
		return
	}

	h.queue = make(worker.Queue)
	h.w.Queue = h.queue

	h.w.Start(ctx)
}

// Emit dispatches a new event.
// Parameter event must be one of:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *eventListener) Emit(event any) {
	h.queue <- event
}

// Subscribe subscribes new event handler to a reflect.Type.
// Parameter eventType must be one of:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *eventListener) Subscribe(eventType reflect.Type, fn func(e any)) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.jobs[eventType]

	if !ok {
		h.jobs[eventType] = []func(e any){fn}
	} else {
		h.jobs[eventType] = append(h.jobs[eventType], fn)
	}
}
