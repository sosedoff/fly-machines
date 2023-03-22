.PHONY: test lint

all: test lint

test:
	go test -cover -race ./...

lint:
	golangci-lint run
