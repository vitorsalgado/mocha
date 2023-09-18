package test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/dzhttp"
	"github.com/vitorsalgado/mocha/v3/dzhttp/httpval"
	"github.com/vitorsalgado/mocha/v3/matcher"
)

func TestCallbacks(t *testing.T) {
	t.Run("should call registered post action", func(t *testing.T) {
		spy := false
		m := dzhttp.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(dzhttp.Post(matcher.URLPath("/test")).
			Callback(func(input *dzhttp.CallbackInput) error {
				require.NotNil(t, input)
				require.NotNil(t, input.App)
				require.NotNil(t, input.URL)
				require.NotNil(t, input.Mock)
				require.NotNil(t, input.RawRequest)
				require.NotNil(t, input.MockedResponse)
				require.NotNil(t, input.ParsedBody)
				require.Equal(t, http.MethodPost, input.RawRequest.Method)
				require.Equal(t, "/test", input.URL.Path)
				require.Equal(t, "hi", input.ParsedBody.(string))

				spy = true

				return nil
			}).
			Reply(dzhttp.OK()))

		res, err := http.Post(fmt.Sprintf("%s/test", m.URL()), httpval.MIMETextPlain, strings.NewReader("hi"))
		require.NoError(t, err)

		defer res.Body.Close()

		require.True(t, scope.AssertCalled(t))
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, spy)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		var callbackErrReceiver error

		callbackErr := errors.New("failed to run post action")
		m := dzhttp.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(dzhttp.Get(matcher.URLPath("/test")).
			Callback(func(input *dzhttp.CallbackInput) error {
				require.NotNil(t, input)
				callbackErrReceiver = callbackErr
				return callbackErr
			}).
			Reply(dzhttp.OK()))

		res, err := http.Get(fmt.Sprintf("%s/test", m.URL()))
		require.NoError(t, err)

		defer res.Body.Close()

		require.True(t, scope.AssertCalled(t))
		require.NotNil(t, callbackErrReceiver)
		require.Equal(t, callbackErr, callbackErrReceiver)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
