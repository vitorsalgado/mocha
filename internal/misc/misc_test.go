package misc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringify(t *testing.T) {
	str := Stringify("test")
	assert.Equal(t, "test", str)

	str = Stringify(true)
	assert.Equal(t, "true", str)

	str = Stringify(10.01)
	assert.Equal(t, "10.01", str)

	var a any
	str = Stringify(a)
	assert.Equal(t, "<value omitted: type=not_defined>", str)

	st := struct {
		name string
	}{name: "hello"}
	str = Stringify(st)
	assert.Equal(t, "<value omitted: type=not_defined>", str)
}
