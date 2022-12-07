package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/reply"
)

func TestPriority(t *testing.T) {
	m := mocha.New(t)
	m.MustStart()

	defer m.Close()

	one := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(3).
		Reply(reply.OK()))
	two := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(1).
		Reply(reply.BadRequest()))
	three := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(100).
		Reply(reply.Created()))

	res, err := testutil.Get(m.URL() + "/test").Do()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.False(t, one.Called())
	assert.True(t, two.Called())
	assert.False(t, three.Called())
}
