package mocha

import (
	"errors"
	"math/rand"
	"net/http"
	"sync"
)

var _ Reply = (*RandomReply)(nil)
var _randomMu sync.Mutex

// RandomReply configures a Reply that serves random HTTP responses.
type RandomReply struct {
	replies []Reply
	random  *rand.Rand
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

func (rep *RandomReply) beforeBuild(_ *Mocha) error {
	size := len(rep.replies)
	if size == 0 {
		return errors.New("reply_random: you need to set at least one response when using random reply")
	}

	return nil
}

// Build builds a response stub randomly based on previously added Reply implementations.
func (rep *RandomReply) Build(w http.ResponseWriter, r *RequestValues) (*Stub, error) {
	_randomMu.Lock()
	defer _randomMu.Unlock()

	var index int
	if rep.random == nil {
		index = r.App.random.Intn(len(rep.replies))
	} else {
		index = rep.random.Intn(len(rep.replies))
	}

	reply := rep.replies[index]

	return reply.Build(w, r)
}
