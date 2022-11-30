package hooks

import (
	"context"
	"reflect"
	"sync"
)

// Hook Types.
var (
	HookOnRequest           = reflect.TypeOf(&OnRequest{})
	HookOnRequestMatched    = reflect.TypeOf(&OnRequestMatch{})
	HookOnRequestNotMatched = reflect.TypeOf(&OnRequestNotMatched{})
	HookOnError             = reflect.TypeOf(&OnError{})
)

type hook reflect.Type
type queue chan any

type Hooks struct {
	worker *worker
	queue  queue
	hooks  map[hook][]func(e any)
	mu     sync.Mutex
}

func New() *Hooks {
	h := &Hooks{}
	h.hooks = map[hook][]func(e any){}
	h.worker = &worker{hooks: h.hooks}

	return h
}

// Start starts background event listener.
func (h *Hooks) Start(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.worker.started.Load() {
		return
	}

	h.queue = make(queue)
	h.worker.queue = h.queue

	h.worker.Start(ctx)
}

// Emit dispatches a new event.
// Parameter event must be one of:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *Hooks) Emit(event any) {
	h.queue <- event
}

// Subscribe subscribes new event handler to a reflect.Type.
// Parameter eventType must be one of:
// - OnRequest
// - OnRequestMatch
// - OnRequestNotMatched
// - OnError
func (h *Hooks) Subscribe(eventType reflect.Type, fn func(e any)) {
	h.mu.Lock()
	defer h.mu.Unlock()

	_, ok := h.hooks[eventType]

	if !ok {
		h.hooks[eventType] = []func(e any){fn}
	} else {
		h.hooks[eventType] = append(h.hooks[eventType], fn)
	}
}
