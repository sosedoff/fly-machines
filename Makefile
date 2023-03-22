.PHONY: setup test lint

all: test lint

setup:
	go mod download

test:
	go test -cover -race ./...

lint:
	golangci-lint run
