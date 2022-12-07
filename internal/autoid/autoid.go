package autoid

import (
	"sync"
)

var id int32
var mu sync.Mutex

func Next() int {
	mu.Lock()
	defer mu.Unlock()

	id++

	return int(id)
}
