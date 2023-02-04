package reply

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnErrorWhenSequenceDoesNotContainReplies(t *testing.T) {
	m := &mmock{}
	m.On("Hits").Return(0)

	res, err := Seq().Build(nil, m, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
