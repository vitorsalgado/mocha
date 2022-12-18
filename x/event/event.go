package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/vitorsalgado/mocha/v3/internal/worker"
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
	Request   *EvtReq
	StartedAt time.Time
}

// OnRequestMatch event is triggered when a mock is found for a request.
type OnRequestMatch struct {
	Request            *EvtReq
	ResponseDefinition EvtRes
	Mock               EvtMk
	Elapsed            time.Duration
}

// OnRequestNotMatched event is triggered when no mocks are found for a request.
type OnRequestNotMatched struct {
	Request *EvtReq
	Result  EvtResult
}

// OnError event is triggered when an error occurs during request matching.
type OnError struct {
	Request *EvtReq
	Err     error
}

type Listener struct {
	w     *worker.Worker
	queue worker.Queue
	jobs  map[worker.JobType][]func(e any)
	mu    sync.Mutex
}

func New() *Listener {
	h := &Listener{}

	h.jobs = map[worker.JobType][]func(e any){}
	h.jobs[EventOnRequest] = make([]func(e any), 0)
	h.jobs[EventOnRequestMatched] = make([]func(e any), 0)
	h.jobs[EventOnRequestNotMatched] = make([]func(e any), 0)
	h.jobs[EventOnError] = make([]func(e any), 0)

	h.w = &worker.Worker{Jobs: h.jobs}

	return h
}

// StartListening starts background event listener.
func (h *Listener) StartListening(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

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
func (h *Listener) Emit(event any) {
	h.queue <- event
}

// Subscribe subscribes new event handler to a reflect.Type.
// Parameter eventType must be one of:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *Listener) Subscribe(eventType reflect.Type, fn func(e any)) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.jobs[eventType]

	if !ok {
		return fmt.Errorf("unknown event type %s", eventType.Name())
	}

	h.jobs[eventType] = append(h.jobs[eventType], fn)

	return nil
}
