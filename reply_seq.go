package mocha

import (
	"fmt"
	"net/http"
	"sync"
)

var _ Reply = (*SequentialReply)(nil)

// SequentialReply configures a sequence of replies to be used after a Mock is matched to a http.Request.
type SequentialReply struct {
	replyAfterSequenceEnded Reply
	replies                 []Reply
	hits                    int
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

func (r *SequentialReply) beforeBuild(_ *Mocha) error {
	size := len(r.replies)
	if size == 0 {
		return fmt.Errorf("reply_sequence: you need to set at least one response when using multiple sequential reply")
	}

	return nil
}

// Build builds a new response based on the current Mock call sequence.
// When the sequence is over, it will return an error or a previously configured reply for this scenario.
func (r *SequentialReply) Build(w http.ResponseWriter, req *RequestValues) (*Stub, error) {
	var size = len(r.replies)
	var reply Reply
	var hits = r.curHits()

	if hits < size && hits >= 0 {
		reply = r.replies[hits]
	}

	if reply == nil {
		if r.replyAfterSequenceEnded != nil {
			return r.replyAfterSequenceEnded.Build(w, req)
		}

		return nil,
			fmt.Errorf(
				"reply_sequence: unable to obtain a response, and no default response was set. request number=%d, sequence size=%d",
				hits,
				size)
	}

	r.updateHits()

	return reply.Build(w, req)
}

func (r *SequentialReply) curHits() int {
	r.rwMutex.RLock()
	defer r.rwMutex.RUnlock()

	return r.hits
}

func (r *SequentialReply) updateHits() {
	r.rwMutex.Lock()
	defer r.rwMutex.Unlock()

	r.hits++
}
