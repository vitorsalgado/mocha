.PHONY: init
init:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/gogo/protobuf/protoc-gen-gofast@latest
	go install github.com/cosmtrek/air@latest

.PHONY: download
download:
	@go mod tidy
	@go mod download

.PHONY: install
install: download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
