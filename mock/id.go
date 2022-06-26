package mock

import "sync/atomic"

type autoID struct {
	id int32
}

func (i *autoID) Next() int {
	atomic.AddInt32(&i.id, 1)
	return int(i.id)
}
