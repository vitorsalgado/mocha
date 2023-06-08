package httpd

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

var _ Reply = (*SequentialReply)(nil)

// SequentialReply configures a sequence of replies to be used after a Mock is matched to a http.Request.
type SequentialReply struct {
	replyAfterSequenceEnded Reply
	replies                 []Reply
	hits                    int32
	rwMutex                 sync.RWMutex
}

// Seq creates a new SequentialReply.
func Seq(reply ...Reply) *SequentialReply {
	return &SequentialReply{replies: reply}
}

// OnSequenceEnded sets a response to be used once the sequence is over.
func (r *SequentialReply) OnSequenceEnded(reply Reply) *SequentialReply {
	r.replyAfterSequenceEnded = reply
	return r
}

// Add adds a new response to the sequence.
func (r *SequentialReply) Add(reply ...Reply) *SequentialReply {
	r.replies = append(r.replies, reply...)
	return r
}

func (r *SequentialReply) beforeBuild(_ *HTTPMockApp) error {
	size := len(r.replies)
	if size == 0 {
		return fmt.Errorf("reply_sequence: you need to set at least one response when using multiple sequential reply")
	}

	return nil
}

// Build builds a new response based on the current Mock call sequence.
// When the sequence is over, it will return an error or a previously configured reply for this scenario.
func (r *SequentialReply) Build(w http.ResponseWriter, req *RequestValues) (*Stub, error) {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	var size = int32(len(r.replies))
	var reply Reply

	if r.hits < size && r.hits >= 0 {
		reply = r.replies[r.hits]
	}

	if reply == nil {
		if r.replyAfterSequenceEnded != nil {
			return r.replyAfterSequenceEnded.Build(w, req)
		}

		return nil,
			fmt.Errorf(
				"reply_sequence: unable to obtain a response, and no default response was set. request number=%d, sequence size=%d",
				r.hits,
				size)
	}

	atomic.AddInt32(&r.hits, 1)

	return reply.Build(w, req)
}

func (r *SequentialReply) totalHits() int {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	return int(r.hits)
}
