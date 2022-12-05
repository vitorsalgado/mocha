package reply

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomReplies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	statuses := []int{
		http.StatusOK, http.StatusInternalServerError, http.StatusCreated, http.StatusBadRequest}

	for i := 0; i < 5000; i++ {
		res, err := Rand(
			BadRequest(),
			OK(),
			Created(),
			InternalServerError(),
		).Build(nil, req)

		contains := false
		for _, status := range statuses {
			if status == res.Status {
				contains = true
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, contains)
	}
}

func TestShouldReturnErrorWhenRandomDoesNotContainReplies(t *testing.T) {
	res, err := Rand().Build(nil, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}
