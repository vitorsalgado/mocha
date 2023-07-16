package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEachSlice(t *testing.T) {
	s := []string{"test", "tested", "testing"}
	r, err := Each(HasPrefix("test")).Match(s)

	require.NoError(t, err)
	require.True(t, r.Pass)
}

func TestEachSliceNotPass(t *testing.T) {
	s := []string{"test", "tested", "dev"}
	r, err := Each(HasPrefix("test")).Match(s)

	require.NoError(t, err)
	require.False(t, r.Pass)
}

func TestEachMap(t *testing.T) {
	m := map[string]string{"1": "test", "2": "testing", "3": "tested"}
	r, err := Each(HasPrefix("test")).Match(m)

	require.NoError(t, err)
	require.True(t, r.Pass)
}

func TestEachMapNotPass(t *testing.T) {
	m := map[string]string{"1": "dev", "2": "testing", "3": "tested"}
	r, err := Each(HasPrefix("test")).Match(m)

	require.NoError(t, err)
	require.False(t, r.Pass)
}

func TestEachInvalidValueType(t *testing.T) {
	_, err := Each(Anything()).Match("hello")
	require.Error(t, err)
}
