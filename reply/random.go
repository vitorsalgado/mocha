package reply

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/vitorsalgado/mocha/v3/types"
)

var _ Reply = (*RandomReply)(nil)

// RandomReply configures a Reply that serves random HTTP responses.
type RandomReply struct {
	replies []Reply
	random  *rand.Rand
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
	r.random = random

	return r
}

// Add adds a new Reply to the random list.
func (r *RandomReply) Add(reply ...Reply) *RandomReply {
	r.replies = append(r.replies, reply...)
	return r
}

func (r *RandomReply) Prepare() error {
	size := len(r.replies)
	if size == 0 {
		return fmt.Errorf("you need to set at least one response when using random reply")
	}

	return nil
}

func (r *RandomReply) Raw() types.RawValue {
	replies := make([]any, len(r.replies))
	for i, rr := range r.replies {
		if rr, ok := rr.(types.Persist); ok {
			replies[i] = rr.Raw().Arguments()
		}
	}

	return types.RawValue{"response_random", map[string]any{
		"responses": replies,
	}}
}

// Build builds a response stub randomly based on previously added Reply implementations.
func (r *RandomReply) Build(w http.ResponseWriter, req *types.RequestValues) (*Stub, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var index int
	if r.random == nil {
		index = rand.Intn(len(r.replies)-1) + 0
	} else {
		index = r.random.Intn(len(r.replies)-1) + 0
	}

	reply := r.replies[index]

	return reply.Build(w, req)
}
