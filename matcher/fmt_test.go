package matcher

import (
	"fmt"
	"testing"
)

func TestIndent(t *testing.T) {
	str := "hello\nworld"
	in := indent(str)

	fmt.Println(in)
}
