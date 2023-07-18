package test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/vitorsalgado/mocha/v3/dzhttp"
	. "github.com/vitorsalgado/mocha/v3/matcher"
)

func TestHTTPMethods(t *testing.T) {
	client := &http.Client{}

	m := NewAPI().CloseWithT(t)
	m.MustStart()

	path1 := "/test/1"
	path2 := "/test/2"
	path3 := "/test/3"

	testCases := []struct {
		method string
		path   string
		mock   *HTTPMockBuilder
		status int
	}{
		{http.MethodGet, path1, Getf(path1), 200},
		{http.MethodGet, path2, Get(URLPathf(path2)), 400},
		{http.MethodGet, path3, Request().Method(http.MethodGet).URLPathf(path3), 500},
		{http.MethodGet, path1, m.Getf(path1), 207},
		{http.MethodGet, path2, m.Get(URLPathf(path2)), 407},

		{http.MethodPost, path1, Postf(path1), 201},
		{http.MethodPost, path2, Post(URLPathf(path2)), 401},
		{http.MethodPost, path3, Request().Method(http.MethodPost).URLPathf(path3), 501},
		{http.MethodPost, path1, m.Postf(path1), 208},
		{http.MethodPost, path2, m.Post(URLPathf(path2)), 408},

		{http.MethodPut, path1, Putf(path1), 202},
		{http.MethodPut, path2, Put(URLPathf(path2)), 402},
		{http.MethodPut, path3, Request().Method(http.MethodPut).URLPathf(path3), 502},
		{http.MethodPut, path1, m.Putf(path1), 209},
		{http.MethodPut, path2, m.Put(URLPathf(path2)), 409},

		{http.MethodPatch, path1, Patchf(path1), 203},
		{http.MethodPatch, path2, Patch(URLPathf(path2)), 403},
		{http.MethodPatch, path3, Request().Method(http.MethodPatch).URLPathf(path3), 503},
		{http.MethodPatch, path1, m.Patchf(path1), 210},
		{http.MethodPatch, path2, m.Patch(URLPathf(path2)), 410},

		{http.MethodDelete, path1, Deletef(path1), 204},
		{http.MethodDelete, path2, Delete(URLPathf(path2)), 404},
		{http.MethodDelete, path3, Request().Method(http.MethodDelete).URLPathf(path3), 504},
		{http.MethodDelete, path1, m.Deletef(path1), 211},
		{http.MethodDelete, path2, m.Delete(URLPathf(path2)), 411},

		{http.MethodHead, path1, Headf(path1), 205},
		{http.MethodHead, path2, Head(URLPathf(path2)), 405},
		{http.MethodHead, path3, Request().Method(http.MethodHead).URLPathf(path3), 505},
		{http.MethodHead, path1, m.Headf(path1), 211},
		{http.MethodHead, path2, m.Head(URLPathf(path2)), 411},

		{http.MethodConnect, path1, Request().Method(http.MethodConnect).URLPathf(path1), 206},
		{http.MethodOptions, path2, Request().Method(http.MethodOptions).URLPathf(path2), 406},
		{http.MethodTrace, path3, Request().Method(http.MethodTrace).URLPathf(path3), 506},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			m.MustMock(tc.mock.Reply(Status(tc.status)))

			req, _ := http.NewRequest(tc.method, m.URL(tc.path), nil)
			res, err := client.Do(req)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, tc.status, res.StatusCode)

			other := http.MethodGet
			if tc.method == http.MethodGet {
				other = http.MethodPost
			}

			req, _ = http.NewRequest(other, m.URL(tc.path), nil)
			res, err = client.Do(req)

			require.NoError(t, err)
			require.NoError(t, res.Body.Close())
			require.Equal(t, StatusNoMatch, res.StatusCode)

			m.Clean()
		})
	}
}
