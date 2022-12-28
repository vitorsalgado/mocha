package mocha

import (
	"fmt"
	"net/http"
	"sync"
)

var _ Reply = (*SequentialReply)(nil)

// SequentialReply configures a sequence of replies to be used after a mock.Mock is matched to a http.Request.
type SequentialReply struct {
	replyAfterSeqEnded Reply
	replies            []Reply
	hits               int
	mu                 sync.Mutex
}

// Seq creates a new SequentialReply.
func Seq(reply ...Reply) *SequentialReply {
	return &SequentialReply{replies: reply}
}

// AfterSequenceEnded sets a response to be used once the sequence is over.
func (r *SequentialReply) AfterSequenceEnded(reply Reply) *SequentialReply {
	r.replyAfterSeqEnded = reply
	return r
}

// Add adds a new response to the sequence.
func (r *SequentialReply) Add(reply ...Reply) *SequentialReply {
	r.replies = append(r.replies, reply...)
	return r
}

func (r *SequentialReply) Pre() error {
	size := len(r.replies)
	if size == 0 {
		return fmt.Errorf("you need to set at least one response when using multiple response builder")
	}

	return nil
}

// Build builds a new response based on current mock.Mock call sequence.
// When the sequence is over, it will return an error or a previously configured reply for this scenario.
func (r *SequentialReply) Build(w http.ResponseWriter, req *RequestValues) (*Stub, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var size = len(r.replies)
	var reply Reply

	if r.hits < size && r.hits >= 0 {
		reply = r.replies[r.hits]
	}

	if reply == nil {
		if r.replyAfterSeqEnded != nil {
			return r.replyAfterSeqEnded.Build(w, req)
		}

		return nil,
			fmt.Errorf(
				"unable to obtain a response and no default response was set. request number: %d - sequence size: %d",
				r.hits,
				size)
	}

	r.hits++

	return reply.Build(w, req)
}
