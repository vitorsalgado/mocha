package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifier_Helper(t *testing.T) {
	n := StdoutNotifier()
	assert.NotPanics(t, n.Helper)
}

func TestNotifier_FailNow(t *testing.T) {
	n := StdoutNotifier()
	assert.NotPanics(t, n.FailNow)
}

func TestNotifier_Errorf(t *testing.T) {
	n := StdoutNotifier()
	assert.NotPanics(t, func() {
		n.Errorf("test %s", "hello")
	})
}
