package internal

import "sync/atomic"

type ID struct {
	id int32
}

func (i *ID) Next() int32 {
	atomic.AddInt32(&i.id, 1)
	return i.id
}
