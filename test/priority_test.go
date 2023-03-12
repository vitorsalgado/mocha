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

func TestPriority_DefaultIsZero(t *testing.T) {
	httpClient := &http.Client{}
	m := mocha.New(mocha.Setup().
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
	m := mocha.New(mocha.Setup().
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
