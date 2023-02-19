package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestSequenceReplies(t *testing.T) {
	m := mocha.New()
	m.MustStart()
	m.MustMock(mocha.Get(URLPath("/test")).
		Reply(mocha.Seq().
			Add(mocha.Unauthorized(), mocha.OK())))

	defer m.Close()

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusTeapot, res.StatusCode)
	}
}

func TestSequenceRepliesOnSequenceEndsSet(t *testing.T) {
	m := mocha.New()
	m.MustStart()
	m.MustMock(mocha.Get(URLPath("/test")).
		Reply(mocha.Seq().
			Add(mocha.Unauthorized(), mocha.OK()).
			OnSequenceEnded(mocha.Created())))

	defer m.Close()

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	for i := 0; i < 2; i++ {
		res, err = http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)
	}
}
