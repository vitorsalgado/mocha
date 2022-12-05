package worker

import (
	"context"
	"log"
	"reflect"
	"sync/atomic"
)

type JobType reflect.Type
type Queue chan any

type Worker struct {
	Started atomic.Bool
	Queue   Queue
	Jobs    map[JobType][]func(e any)
}

func (w *Worker) Start(ctx context.Context) {
	w.Started.Store(true)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()

		for {
			select {
			case event, ok := <-w.Queue:
				if !ok {
					return
				}

				t := reflect.TypeOf(event)
				fns, ok := w.Jobs[t]
				if !ok {
					continue
				}

				for _, fn := range fns {
					fn(event)
				}
			case <-ctx.Done():
				close(w.Queue)
				w.Started.Store(false)
				return
			}
		}
	}()
}
