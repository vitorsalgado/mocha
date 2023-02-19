package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/internal/testutil"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestPriority(t *testing.T) {
	m := mocha.New()
	m.MustStart()

	defer m.Close()

	one := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(3).
		Reply(mocha.OK()))

	two := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(1).
		Reply(mocha.BadRequest()))

	three := m.MustMock(mocha.Get(matcher.URLPath("/test")).
		Priority(100).
		Reply(mocha.Created()))

	for i := 0; i < 5; i++ {
		res, err := testutil.Get(m.URL() + "/test").Do()

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.False(t, one.HasBeenCalled())
		assert.True(t, two.HasBeenCalled())
		assert.False(t, three.HasBeenCalled())
	}
}
