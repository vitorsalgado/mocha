package reply

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

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
			if status == res.StatusCode {
				contains = true
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, contains)
	}
}

func TestShouldReturnErrorWhenRandomDoesNotContainReplies(t *testing.T) {
	assert.Error(t, Rand().Prepare())
}

func TestRandWithCustom(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	statuses := []int{
		http.StatusOK, http.StatusInternalServerError, http.StatusCreated, http.StatusBadRequest}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 5000; i++ {
		res, err := RandWithCustom(
			r,
			BadRequest(),
			OK(),
			Created(),
			InternalServerError(),
		).Build(nil, req)

		contains := false
		for _, status := range statuses {
			if status == res.StatusCode {
				contains = true
				break
			}
		}

		assert.Nil(t, err)
		assert.True(t, contains)
	}
}
