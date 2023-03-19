package test

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vitorsalgado/mocha/v3"
	"github.com/vitorsalgado/mocha/v3/matcher"
	"github.com/vitorsalgado/mocha/v3/misc"
)

func TestCallbacks(t *testing.T) {
	t.Run("should call registered post action", func(t *testing.T) {
		spy := false
		m := mocha.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mocha.Post(matcher.URLPath("/test")).
			Callback(func(input *mocha.CallbackInput) error {
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
			Reply(mocha.OK()))

		res, err := http.Post(fmt.Sprintf("%s/test", m.URL()), misc.MIMETextPlain, strings.NewReader("hi"))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.Equal(t, http.StatusOK, res.StatusCode)
		require.True(t, spy)
	})

	t.Run("should not be affected by errors on registered post actions", func(t *testing.T) {
		var callbackErrReceiver error

		callbackErr := errors.New("failed to run post action")
		m := mocha.NewAPI()
		m.MustStart()

		defer m.Close()

		scope := m.MustMock(mocha.Get(matcher.URLPath("/test")).
			Callback(func(input *mocha.CallbackInput) error {
				require.NotNil(t, input)
				callbackErrReceiver = callbackErr
				return callbackErr
			}).
			Reply(mocha.OK()))

		res, err := http.Get(fmt.Sprintf("%s/test", m.URL()))

		require.NoError(t, err)
		require.True(t, scope.AssertCalled(t))
		require.NotNil(t, callbackErrReceiver)
		require.Equal(t, callbackErr, callbackErrReceiver)
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
