package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestSequenceReplies(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.Get(URLPath("/test")).
		Reply(dzhttp.Seq().
			Add(dzhttp.Unauthorized(), dzhttp.OK())))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusTeapot, res.StatusCode)
	}
}

func TestSequenceRepliesOnSequenceEndsSet(t *testing.T) {
	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.Get(URLPath("/test")).
		Reply(dzhttp.Seq().
			Add(dzhttp.Unauthorized(), dzhttp.OK()).
			OnSequenceEnded(dzhttp.Created())))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	for i := 0; i < 2; i++ {
		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusCreated, res.StatusCode)
	}
}

func TestSequence_SetupFromFile(t *testing.T) {
	type response struct {
		Ok   bool    `json:"ok,omitempty"`
		Type string  `json:"type,omitempty"`
		Num  float64 `json:"num,omitempty"`
	}

	m := dzhttp.NewAPI()
	m.MustStart()

	defer m.Close()

	m.MustMock(dzhttp.FromFile("testdata/response_sequence/seq_01.yaml"))

	req, _ := http.NewRequest(http.MethodGet, m.URL()+"/test", nil)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())

	r := new(response)
	err = json.Unmarshal(b, r)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, true, r.Ok)
	require.Equal(t, "test", r.Type)
	require.EqualValues(t, 10, r.Num)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusCreated, res.StatusCode)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	for i := 0; i < 3; i++ {
		res, err = http.DefaultClient.Do(req)
		require.NoError(t, err)
		require.NoError(t, res.Body.Close())
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	}
}
