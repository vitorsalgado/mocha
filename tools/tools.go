//go:build tools
// +build tools

package tools

import (
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "golang.org/x/vuln"
)
