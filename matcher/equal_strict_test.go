package matcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStrictEqual(t *testing.T) {
	job := "jobless"
	timestamp, err := time.Parse("2006-01-02T15:04:05.000Z", "2022-12-31T00:30:25.010Z")
	require.NoError(t, err)

	type model struct {
		name      string
		year      int
		value     float64
		timestamp time.Time
		active    bool
		job       *string
	}

	testCases := []struct {
		name           string
		value          any
		valueToCompare any
		result         bool
	}{
		{"string with nil", "test", nil, false},
		{"nil", nil, nil, true},
		{"bool (true)", true, true, true},
		{"bool (false)", true, false, false},
		{"byte arrays", []byte("test"), []byte("test"), true},
		{"array", []string{"dev", "test", "hello world"}, []string{"dev", "test", "hello world"}, true},
		{"array (diff order)", []string{"dev", "test", "hello world"}, []string{"test", "dev", "hello world"}, false},
		{"map", map[string]any{"active": true, "name": "test"}, map[string]any{"active": true, "name": "test"}, true},
		{"map (diff order)", map[string]any{"active": true, "name": "test"}, map[string]any{"name": "test", "active": true}, true},
		{"struct", model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, true},
		{"struct (diff order)", model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, model{timestamp: timestamp, active: true, year: 2022, name: "test", value: 100.255, job: nil}, true},
		{"struct (false)", model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, model{value: 100.255, timestamp: timestamp, active: true, job: nil}, false},
		{"struct (diff one field)", model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: &job}, false},
		{"struct (pinter and non-pointer", &model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, false},
		{"struct (two pointers)", &model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, &model{name: "test", year: 2022, value: 100.255, timestamp: timestamp, active: true, job: nil}, true},
		{"diff types", []any{"dev", "test"}, []string{"dev", "test"}, false},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := StrictEqual(tc.value).Match(tc.valueToCompare)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)

			res, err = Eqs(tc.valueToCompare).Match(tc.value)

			require.NoError(t, err)
			require.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestStrictEqualf(t *testing.T) {
	result, err := StrictEqualf("hello %s %d", "world", 10).Match("hello world 10")
	require.NoError(t, err)
	require.True(t, result.Pass)
}
