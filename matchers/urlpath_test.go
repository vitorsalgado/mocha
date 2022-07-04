package matchers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlPath(t *testing.T) {
	u, _ := url.Parse("http://localhost:8080/test/hello")
	result, err := URLPath("/test/hello").Matches(*u, Args{})

	assert.Nil(t, err)
	assert.True(t, result)
}
