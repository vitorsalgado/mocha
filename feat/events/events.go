package events

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"
)

// Event data transfer objects
type (
	Request struct {
		Method     string
		Path       string
		RequestURI string
		Host       string
		Header     http.Header
	}

	Response struct {
		Status int
		Header http.Header
	}

	Mock struct {
		ID   int
		Name string
	}

	ResultDetail struct {
		Name        string
		Target      string
		Description string
	}

	Result struct {
		ClosestMatch *Mock
		Details      []ResultDetail
	}
)

// Events
type (
	OnRequest struct {
		Request   Request
		StartedAt time.Time
	}

	OnRequestMatch struct {
		Request            Request
		ResponseDefinition Response
		Mock               Mock
		Elapsed            time.Duration
	}

	OnRequestNotMatched struct {
		Request Request
		Result  Result
	}

	OnError struct {
		Request Request
		Err     error
	}
)

type (
	Events interface {
		OnRequest(OnRequest)
		OnRequestMatch(OnRequestMatch)
		OnRequestNotMatched(OnRequestNotMatched)
		OnError(OnError)
	}

	Emitter struct {
		ctx      context.Context
		events   []Events
		listener chan any
		mu       sync.Mutex
	}
)

func NewEmitter(ctx context.Context) *Emitter {
	return &Emitter{ctx: ctx, events: make([]Events, 0), listener: make(chan any)}
}

func (h *Emitter) Emit(data any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	go func(data any, listener chan any) { listener <- data }(data, h.listener)
}

func (h *Emitter) Subscribe(evt Events) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, evt)
}

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
						panic(fmt.Errorf("event type %s is invalid", reflect.TypeOf(data).Name()))
					}
				}
			}
		}
	}()
}
