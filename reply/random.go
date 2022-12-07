package reply

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
)

var _ Reply = (*RandomReply)(nil)

// RandomReply configures a Reply that serves random HTTP responses.
type RandomReply struct {
	replies []Reply
	r       *rand.Rand
	mu      sync.Mutex
}

// Rand inits a new RandomReply.
func Rand(reply ...Reply) *RandomReply {
	return &RandomReply{
		replies: reply,
	}
}

// RandWithCustom creates a new RandomReply with a custom *rand.Rand.
func RandWithCustom(random *rand.Rand, reply ...Reply) *RandomReply {
	r := Rand(reply...)
	r.r = random

	return r
}

// Add adds a new Reply to the random list.
func (mr *RandomReply) Add(reply ...Reply) *RandomReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

// Build builds a response stub randomly based on previously added Reply implementations.
func (mr *RandomReply) Build(w http.ResponseWriter, r *http.Request) (*Response, error) {
	size := len(mr.replies)
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using random reply")
	}

	mr.mu.Lock()
	defer mr.mu.Unlock()

	var index int
	if mr.r == nil {
		index = rand.Intn(len(mr.replies)-1) + 0
	} else {
		index = mr.r.Intn(len(mr.replies)-1) + 0
	}

	reply := mr.replies[index]

	return reply.Build(w, r)
}
