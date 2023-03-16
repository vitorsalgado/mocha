package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestPriority(t *testing.T) {
	m := New()
	m.MustStart()

	defer m.Close()

	one := m.MustMock(Get(URLPath("/test")).
		Priority(3).
		Reply(OK()))

	two := m.MustMock(Get(URLPath("/test")).
		Priority(1).
		Reply(BadRequest()))

	three := m.MustMock(Get(URLPath("/test")).
		Priority(100).
		Reply(Created()))

	for i := 0; i < 5; i++ {
		res, err := http.Get(m.URL() + "/test")

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.False(t, one.HasBeenCalled())
		assert.True(t, two.HasBeenCalled())
		assert.False(t, three.HasBeenCalled())
	}
}

func TestPriority_DefaultIsZero(t *testing.T) {
	httpClient := &http.Client{}
	m := New(Setup().
		MockFilePatterns(
			"testdata/priority/default_is_zero/*.json",
			"testdata/priority/default_is_zero/*.yaml"))
	m.MustStart()

	defer m.Close()

	testCases := []struct {
		status int
		url    string
	}{
		{http.StatusNoContent, "/test"},
		{http.StatusNoContent, "/test/hello"},
		{http.StatusNoContent, "/test/hello/world"},
		{http.StatusNotFound, "/hello"},
		{http.StatusNotFound, "/TEST"},
		{http.StatusNotFound, "/TEST/hello"},
		{http.StatusNotFound, "/TEST/HELLO"},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			res, err := httpClient.Get(m.URL() + tc.url)

			require.NoError(t, err)
			require.Equal(t, tc.status, res.StatusCode)
		})
	}
}

func TestPriority_LowestShouldBeServed(t *testing.T) {
	httpClient := &http.Client{}
	m := New(Setup().
		RootDir("testdata/priority").
		MockFilePatterns(
			"lowest_should_be_served/*.json",
			"lowest_should_be_served/*.yaml"))
	m.MustStart()

	defer m.Close()

	testCases := []struct {
		status int
		url    string
	}{
		{http.StatusNoContent, "/test"},
		{http.StatusNoContent, "/test/hello"},
		{http.StatusNoContent, "/test/hello/world"},
		{http.StatusNotFound, "/hello"},
		{http.StatusNotFound, "/TEST"},
		{http.StatusNotFound, "/TEST/hello"},
		{http.StatusNotFound, "/TEST/HELLO"},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			res, err := httpClient.Get(m.URL() + tc.url)

			require.NoError(t, err)
			require.Equal(t, tc.status, res.StatusCode)
		})
	}
}
