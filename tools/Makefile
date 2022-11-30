download:
	@go mod tidy
	@go mod download

tools: download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
