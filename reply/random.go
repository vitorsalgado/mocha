package reply

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/vitorsalgado/mocha/core"
	"github.com/vitorsalgado/mocha/internal/parameters"
)

// RandomReply configures a core.Reply that serves random HTTP responses.
type RandomReply struct {
	replies []core.Reply
}

// Rand inits a new RandomReply.
func Rand() *RandomReply { return &RandomReply{replies: make([]core.Reply, 0)} }

// Add adds a new core.Reply to the random list.
func (mr *RandomReply) Add(reply ...core.Reply) *RandomReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

// Build builds a response stub randomly based on previously added core.Reply implementations.
func (mr *RandomReply) Build(r *http.Request, m *core.Mock, p parameters.Params) (*core.Response, error) {
	size := len(mr.replies)
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using random reply")
	}

	index := rand.Intn(len(mr.replies)-1) + 0
	reply := mr.replies[index]

	return reply.Build(r, m, p)
}
