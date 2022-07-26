package autoid

import "sync/atomic"

var id int32

func Next() int {
	atomic.AddInt32(&id, 1)
	return int(id)
}
