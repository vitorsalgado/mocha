package reply

import (
	"fmt"
	"net/http"

	"github.com/vitorsalgado/mocha/v3/params"
)

// SequentialReply configures a sequence of replies to be used after a mock.Mock is matched to a http.Request.
type SequentialReply struct {
	replyOnNotFound Reply
	replies         []Reply
}

// Seq creates a new SequentialReply.
func Seq() *SequentialReply {
	return &SequentialReply{replies: make([]Reply, 0)}
}

// AfterEnded sets a response to be used once the sequence is over.
func (mr *SequentialReply) AfterEnded(reply Reply) *SequentialReply {
	mr.replyOnNotFound = reply
	return mr
}

// Add adds a new response to the sequence.
func (mr *SequentialReply) Add(reply ...Reply) *SequentialReply {
	mr.replies = append(mr.replies, reply...)
	return mr
}

// Build builds a new response based on current mock.Mock call sequence.
// When the sequence is over, it will return an error or a previously configured reply for this scenario.
func (mr *SequentialReply) Build(r *http.Request, m M, p params.P) (*Response, error) {
	size := len(mr.replies)
	hits := m.Hits()
	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using multiple response builder")
	}

	var reply Reply

	if hits < size {
		reply = mr.replies[hits]
	}

	if reply == nil {
		if mr.replyOnNotFound != nil {
			return mr.replyOnNotFound.Build(r, m, p)
		}

		return nil,
			fmt.Errorf(
				"unable to obtain a response and no default response was set. request number: %d - sequence size: %d",
				hits,
				size)
	}

	return reply.Build(r, m, p)
}
