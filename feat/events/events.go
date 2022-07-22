package events

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Event data transfer objects
type (
	// Request defines information from http.Request to be logged.
	Request struct {
		Method     string
		Path       string
		RequestURI string
		Host       string
		Header     http.Header
	}

	// Response defines HTTP Response information to be logged.
	Response struct {
		Status int
		Header http.Header
	}

	// Mock defines core.Mock information to be logged.
	Mock struct {
		ID   int
		Name string
	}

	// ResultDetail defines matching result details to be logged.
	ResultDetail struct {
		Name        string
		Target      string
		Description string
	}

	// Result defines matching result to be logged.
	Result struct {
		HasClosestMatch bool
		ClosestMatch    Mock
		Details         []ResultDetail
	}
)

// Events
type (
	// OnRequest event is triggered every time a request arrives at the mock handler.
	OnRequest struct {
		Request   Request
		StartedAt time.Time
	}

	// OnRequestMatch event is triggered when a mock is found for a request.
	OnRequestMatch struct {
		Request            Request
		ResponseDefinition Response
		Mock               Mock
		Elapsed            time.Duration
	}

	// OnRequestNotMatched event is triggered when no mocks are found for a request.
	OnRequestNotMatched struct {
		Request Request
		Result  Result
	}

	// OnError event is triggered when an error occurs during request matching.
	OnError struct {
		Request Request
		Err     error
	}
)

type (
	// Events interface defines available event handlers.
	Events interface {
		OnRequest(OnRequest)
		OnRequestMatch(OnRequestMatch)
		OnRequestNotMatched(OnRequestNotMatched)
		OnError(OnError)
	}

	// Emitter implements a event listener, subscriber and emitter.
	Emitter struct {
		ctx      context.Context
		events   []Events
		listener chan any
		mu       sync.Mutex
	}
)

// NewEmitter creates an Emitter instance.
func NewEmitter(ctx context.Context) *Emitter {
	return &Emitter{ctx: ctx, events: make([]Events, 0), listener: make(chan any)}
}

// Emit dispatches a new event.
// Parameter data must be:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *Emitter) Emit(data any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	go func(data any, listener chan any) { listener <- data }(data, h.listener)
}

// Subscribe subscribes new event handlers.
func (h *Emitter) Subscribe(evt Events) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, evt)
}

// Start starts event listener on another go routine.
func (h *Emitter) Start() {
	go func() {
		for {
			select {
			case <-h.ctx.Done():
				return

			case data, ok := <-h.listener:
				if !ok {
					return
				}

				for _, hook := range h.events {
					switch evt := data.(type) {
					case OnRequest:
						hook.OnRequest(evt)
					case OnRequestMatch:
						hook.OnRequestMatch(evt)
					case OnRequestNotMatched:
						hook.OnRequestNotMatched(evt)
					case OnError:
						hook.OnError(evt)

					default:
						log.Printf("event type %s is invalid\n", reflect.TypeOf(data).Name())
					}
				}
			}
		}
	}()
}
