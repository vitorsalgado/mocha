package test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3/matcher"
	mhttp2 "github.com/vitorsalgado/mocha/v3/mhttp"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestCallbacks(t *testing.T) {
	t.Run("should call registered post action", func(t *testing.T) {
		spy := false
		m := mhttp2.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mhttp2.Post(matcher.URLPath("/test")).
			Callback(func(input *mhttp2.CallbackInput) error {
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
			Reply(mhttp2.OK()))

		res, err := http.Post(fmt.Sprintf("%s/test", m.URL()), misc.MIMETextPlain, strings.NewReader("hi"))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, spy)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		var callbackErrReceiver error

		callbackErr := errors.New("failed to run post action")
		m := mhttp2.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mhttp2.Get(matcher.URLPath("/test")).
			Callback(func(input *mhttp2.CallbackInput) error {
				require.NotNil(t, input)
				callbackErrReceiver = callbackErr
				return callbackErr
			}).
			Reply(mhttp2.OK()))

		res, err := http.Get(fmt.Sprintf("%s/test", m.URL()))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.NotNil(t, callbackErrReceiver)
		require.Equal(t, callbackErr, callbackErrReceiver)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
