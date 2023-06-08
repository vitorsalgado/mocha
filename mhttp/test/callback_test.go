package test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/mhttpv"
)

func TestCallbacks(t *testing.T) {
	t.Run("should call registered post action", func(t *testing.T) {
		spy := false
		m := mhttp.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mhttp.Post(matcher.URLPath("/test")).
			Callback(func(input *mhttp.CallbackInput) error {
				require.NotNil(t, input)
				require.NotNil(t, input.App)
				require.NotNil(t, input.URL)
				require.NotNil(t, input.Mock)
				require.NotNil(t, input.RawRequest)
				require.NotNil(t, input.Stub)
				require.NotNil(t, input.ParsedBody)
				require.Equal(t, http.MethodPost, input.RawRequest.Method)
				require.Equal(t, "/test", input.URL.Path)
				require.Equal(t, "hi", input.ParsedBody.(string))

				spy = true

				return nil
			}).
			Reply(mhttp.OK()))

		res, err := http.Post(fmt.Sprintf("%s/test", m.URL()), mhttpv.MIMETextPlain, strings.NewReader("hi"))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, spy)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		var callbackErrReceiver error

		callbackErr := errors.New("failed to run post action")
		m := mhttp.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mhttp.Get(matcher.URLPath("/test")).
			Callback(func(input *mhttp.CallbackInput) error {
				require.NotNil(t, input)
				callbackErrReceiver = callbackErr
				return callbackErr
			}).
			Reply(mhttp.OK()))

		res, err := http.Get(fmt.Sprintf("%s/test", m.URL()))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.NotNil(t, callbackErrReceiver)
		require.Equal(t, callbackErr, callbackErrReceiver)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
