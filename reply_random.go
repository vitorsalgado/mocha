package mocha

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
	random  *rand.Rand
	mu      sync.Mutex
}

// Rand initializes a new RandomReply.
func Rand(reply ...Reply) *RandomReply {
	return &RandomReply{
		replies: reply,
	}
}

// RandWith creates a new RandomReply with a custom *rand.Rand.
func RandWith(random *rand.Rand, reply ...Reply) *RandomReply {
	r := Rand(reply...)
	r.random = random

	return r
}

// Add adds a new Reply to the random list.
func (rep *RandomReply) Add(reply ...Reply) *RandomReply {
	rep.replies = append(rep.replies, reply...)
	return rep
}

func (rep *RandomReply) Pre() error {
	size := len(rep.replies)
	if size == 0 {
		return fmt.Errorf("you need to set at least one response when using random reply")
	}

	return nil
}

// Build builds a response stub randomly based on previously added Reply implementations.
func (rep *RandomReply) Build(w http.ResponseWriter, r *RequestValues) (*Stub, error) {
	rep.mu.Lock()
	defer rep.mu.Unlock()

	var index int
	if rep.random == nil {
		index = rand.Intn(len(rep.replies)-1) + 0
	} else {
		index = rep.random.Intn(len(rep.replies)-1) + 0
	}

	reply := rep.replies[index]

	return reply.Build(w, r)
}
