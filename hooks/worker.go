package hooks

import (
	"context"
	"log"
	"reflect"
	"sync/atomic"
)

type worker struct {
	started atomic.Bool
	queue   queue
	hooks   map[hook][]func(e any)
}

func (w *worker) start(ctx context.Context) {
	w.started.Store(true)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		for {
			select {
			case event, ok := <-w.queue:
				if !ok {
					return
				}

				t := reflect.TypeOf(event)
				fns, ok := w.hooks[t]
				if !ok {
					continue
				}

				for _, fn := range fns {
					fn(event)
				}
			case <-ctx.Done():
				close(w.queue)
				w.started.Store(false)
				return
			}
		}
	}()
}
