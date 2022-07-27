package reply

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/vitorsalgado/mocha/internal/params"
)

// RandomReply configures a Reply that serves random HTTP responses.
type RandomReply struct {
	replies []Reply
}

// Rand inits a new RandomReply.
func Rand() *RandomReply { return &RandomReply{replies: make([]Reply, 0)} }

// Add adds a new Reply to the random list.
func (mr *RandomReply) Add(reply ...Reply) *RandomReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

// Build builds a response stub randomly based on previously added Reply implementations.
func (mr *RandomReply) Build(r *http.Request, m M, p params.P) (*Response, error) {
	size := len(mr.replies)
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using random reply")
	}

	index := rand.Intn(len(mr.replies)-1) + 0
	reply := mr.replies[index]

	return reply.Build(r, m, p)
}
