//go:build tools
// +build tools

package tools

import (
	_ "github.com/cosmtrek/air"
	_ "github.com/gogo/protobuf/protoc-gen-gofast"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
