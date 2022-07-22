package test

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha"
	"github.com/vitorsalgado/mocha/expect"
	"github.com/vitorsalgado/mocha/internal/testutil"
	"github.com/vitorsalgado/mocha/reply"
)

func TestPriority(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	one := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Priority(3).
		Reply(reply.OK()))
	two := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Priority(1).
		Reply(reply.BadRequest()))
	three := m.AddMocks(mocha.Get(expect.URLPath("/test")).
		Priority(100).
		Reply(reply.Created()))

	res, err := testutil.Get(m.URL() + "/test").Do()
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.False(t, one.Called())
	assert.True(t, two.Called())
	assert.False(t, three.Called())
}
