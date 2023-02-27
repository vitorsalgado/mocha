package mocha

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandomReplies(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	rv := &RequestValues{RawRequest: req, URL: req.URL}
	statuses := []int{
		http.StatusOK, http.StatusInternalServerError, http.StatusCreated, http.StatusBadRequest}

	for i := 0; i < 5000; i++ {
		res, err := Rand(
			BadRequest(),
			OK(),
			Created()).
			Add(InternalServerError()).
			Build(nil, rv)

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
	assert.Error(t, Rand().validate(nil))
}

func TestRandWithCustom(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	rv := &RequestValues{RawRequest: req, URL: req.URL}
	statuses := []int{
		http.StatusOK, http.StatusInternalServerError, http.StatusCreated, http.StatusBadRequest}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 5000; i++ {
		res, err := RandWith(
			r,
			BadRequest(),
			OK(),
			Created(),
			InternalServerError(),
		).Build(nil, rv)

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

func TestRandomReplyValidate(t *testing.T) {
	r := Rand()

	require.Error(t, r.validate(nil))

	r.Add(OK())

	require.NoError(t, r.validate(nil))
}
