
.PHONY: all
all: dep build test lint

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: build
build:
	go build ./

.PHONY: test
test:
	go test -race ./...

.PHONY: lint
lint:
	golangci-lint run