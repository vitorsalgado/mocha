package core

import "sync/atomic"

type autoID struct {
	id int32
}

func (i *autoID) next() int {
	atomic.AddInt32(&i.id, 1)
	return int(i.id)
}

var id = autoID{}
