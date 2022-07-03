package reply

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/vitorsalgado/mocha/internal/params"
	"github.com/vitorsalgado/mocha/mock"
)

// RandomReply configures a mock.Reply that serves random HTTP responses.
type RandomReply struct {
	replies []mock.Reply
}

// Rand inits a new RandomReply.
func Rand() *RandomReply { return &RandomReply{replies: make([]mock.Reply, 0)} }

// Add adds a new mock.Reply to the random list.
func (mr *RandomReply) Add(reply ...mock.Reply) *RandomReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

// Build builds a response stub randomly based on previously added mock.Reply implementations.
func (mr *RandomReply) Build(r *http.Request, m *mock.Mock, p params.Params) (*mock.Response, error) {
	size := len(mr.replies)
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using random reply")
	}

	index := rand.Intn(len(mr.replies)-1) + 0
	reply := mr.replies[index]

	return reply.Build(r, m, p)
}
