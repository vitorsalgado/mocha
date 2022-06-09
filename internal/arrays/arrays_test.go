package arrays

import (
	"strings"
	"testing"
)

func TestItShouldMapItemsAndReturnNewTransformedArray(t *testing.T) {
	arr := []string{"1-test", "2-test", "3-test"}
	r1 := All(arr, func(i string) bool { return strings.Contains(i, "-test") })
	r2 := All(arr, func(i string) bool { return strings.Contains(i, "2") })

	if !r1 {
		t.Errorf("expected result to be true but got %t", r1)
	}

	if r2 {
		t.Errorf("expected result to be failse but got %t", r2)
	}
}
