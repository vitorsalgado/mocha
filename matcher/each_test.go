package matcher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEach_Slice(t *testing.T) {
	s := []string{"test", "tested", "testing"}
	r, err := Each(HasPrefix("test")).Match(s)

	assert.NoError(t, err)
	assert.True(t, r.OK)
}

func TestEach_Slice_NotPass(t *testing.T) {
	s := []string{"test", "tested", "dev"}
	r, err := Each(HasPrefix("test")).Match(s)

	assert.NoError(t, err)
	assert.False(t, r.OK)
}

func TestEach_Map(t *testing.T) {
	m := map[string]string{"1": "test", "2": "testing", "3": "tested"}
	r, err := Each(HasPrefix("test")).Match(m)

	assert.NoError(t, err)
	assert.True(t, r.OK)
}

func TestEach_Map_NotPass(t *testing.T) {
	m := map[string]string{"1": "dev", "2": "testing", "3": "tested"}
	r, err := Each(HasPrefix("test")).Match(m)

	assert.NoError(t, err)
	assert.False(t, r.OK)
}
