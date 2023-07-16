package matcher

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToEqualJSON(t *testing.T) {
	testCases := []struct {
		a        any
		b        any
		expected bool
	}{
		{
			map[string]interface{}{"ok": true, "name": "dev"},
			map[string]interface{}{"ok": true, "name": "dev"},
			true,
		},
		{
			map[string]interface{}{"ok": true, "name": "dev"},
			map[string]interface{}{"nok": true, "name": "dev"},
			false,
		},

		{
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			true,
		},
		{
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			`{"name":"no-one","job":{"title":"qa"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			false,
		},
		{
			`{"name":"no-one","job":{"title":"dev","level": 3},"tags":["hi","hello","tchau","bye"],"address":{"number":100,"street":"nowhere","city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			`{"name":"no-one","job":{"level": 3,"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			true,
		},
		{
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.544Z"}`,
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			false,
		},
		{
			`{"name":"no-one","job":{"title":"dev"},"tags":["hello","tchau","bye","hi"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			`{"name":"no-one","job":{"title":"dev"},"tags":["hi","hello","tchau","bye"],"address":{"street":"nowhere","number":100,"city":{"name":"Berlin"}},"active":true,"balance":1250.75,"timestamp":1679084156116,"last_update":"2023-03-17T20:16:21.584Z"}`,
			false,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.FormatInt(int64(i), 10), func(t *testing.T) {
			result, err := EqualJSON(tc.a).Match(tc.b)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result.Pass)
		})
	}
}

func TestEqualJSON_Errors(t *testing.T) {
	ch := make(chan struct{})
	result, err := Eqj(ch).Match(map[string]interface{}{"ok": true, "name": "dev"})

	require.Error(t, err)
	require.False(t, result.Pass)
}
