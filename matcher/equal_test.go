package matcher

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
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
		{"string with float64", "10", float64(10), true},
		{"string with bool", "true", true, false},
		{"float64 with string", float64(10), "10", true},
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
		{"diff types (num)", 10, float64(10), true},
		{"diff types (str)", "test", any("test"), true},
		{"diff types (array)", []any{"dev", "test"}, []string{"dev", "test"}, true},

		{
			"diff types (complex array)",
			[]any{"dev", "test", []any{100, float64(10), int16(2)}},
			[]any{"dev", "test", []any{100, 10, 2}},
			true,
		},
		{
			"diff types (complex array with objects)",
			[]any{"dev", "test", []any{100, float64(10), int16(2), map[string]any{"name": "test", "active": true, "sub": map[string]any{"text": "hello", "num": 100.50}, "level1": []any{10, 20.5}}}},
			[]any{"dev", "test", []any{100, float64(10), int16(2), map[string]any{"name": "test", "active": true, "sub": map[string]any{"text": "hello", "num": 100.50}, "level1": []any{10, 20.5}}}},
			true,
		},
		{
			"diff types (complex array with objects(2))",
			[]any{"dev", "test", []any{100, float64(10), int16(2), map[string]any{"name": "test", "active": true, "sub": map[string]any{"text": "hello", "num": 100.50}, "level1": []any{10, 20.5}}}},
			[]any{"dev", "test", []any{100, 10, 2, map[string]any{"name": "test", "active": true, "sub": map[string]any{"text": "hello", "num": 100.50}, "level1": []float64{10, 20.5}}}},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Equal(tc.value).Match(tc.valueToCompare)

			assert.NoError(t, err)
			assert.Equal(t, tc.result, res.Pass)

			res, err = Equal(tc.valueToCompare).Match(tc.value)

			assert.NoError(t, err)
			assert.Equal(t, tc.result, res.Pass)
		})
	}
}

func TestEqualf(t *testing.T) {
	tcs := []struct {
		name     string
		format   string
		args     []any
		value    any
		expected bool
	}{
		{"", "hello world %s", []any{"dev"}, "hello world dev", true},
		{"", "hello world %d %v", []any{10, true}, "hello world 10 true", true},
		{"", "hello world", []any{"dev", "qa"}, "hello world", false},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := Equalf(tc.format, tc.args...).Match(tc.value)
			require.NoError(t, err)
			require.Equal(t, tc.expected, res.Pass)
		})
	}
}
