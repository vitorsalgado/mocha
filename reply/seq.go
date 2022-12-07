package reply

import (
	"fmt"
	"net/http"
)

var _ Reply = (*SequentialReply)(nil)

// SequentialReply configures a sequence of replies to be used after a mock.Mock is matched to a http.Request.
type SequentialReply struct {
	replyOnNotFound Reply
	replies         []Reply
}

// Seq creates a new SequentialReply.
func Seq(reply ...Reply) *SequentialReply {
	return &SequentialReply{replies: reply}
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

func (mr *SequentialReply) Prepare() error {
	size := len(mr.replies)

	if size == 0 {
		return fmt.Errorf("you need to set at least one response when using multiple response builder")
	}

	return nil
}

// Build builds a new response based on current mock.Mock call sequence.
// When the sequence is over, it will return an error or a previously configured reply for this scenario.
func (mr *SequentialReply) Build(w http.ResponseWriter, r *http.Request) (*Response, error) {
	arg := r.Context().Value(KArg).(*Arg)
	size := len(mr.replies)
	hits := arg.M.Hits

	if size == 0 {
		return nil,
			fmt.Errorf("you need to set at least one response when using multiple response builder")
	}

	var reply Reply

	if hits <= size && hits >= 0 {
		reply = mr.replies[hits-1]
	}

	if reply == nil {
		if mr.replyOnNotFound != nil {
			return mr.replyOnNotFound.Build(w, r)
		}

		return nil,
			fmt.Errorf(
				"unable to obtain a response and no default response was set. request number: %d - sequence size: %d",
				hits,
				size)
	}

	return reply.Build(w, r)
}
