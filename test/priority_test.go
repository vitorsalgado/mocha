package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/expect"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestPriority(t *testing.T) {
	m := mocha.New(t)
	m.Start()

	defer m.Close()

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
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.False(t, one.Called())
	assert.True(t, two.Called())
	assert.False(t, three.Called())
}
