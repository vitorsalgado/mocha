package mocha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifier_Helper(t *testing.T) {
	n := NewConsoleNotifier()
	assert.NotPanics(t, n.Helper)
}

func TestNotifier_FailNow(t *testing.T) {
	n := NewConsoleNotifier()
	assert.NotPanics(t, n.FailNow)
}
