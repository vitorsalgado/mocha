package notifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotifier_Helper(t *testing.T) {
	n := NewConsole()
	assert.NotPanics(t, n.Helper)
}
