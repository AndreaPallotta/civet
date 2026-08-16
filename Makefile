BINARY_NAME=civet
VERSION?=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X github.com/AndreaPallotta/civet/cmd.Version=$(VERSION) -X github.com/AndreaPallotta/civet/cmd.Commit=$(COMMIT) -X github.com/AndreaPallotta/civet/cmd.Date=$(DATE)"

.PHONY: all build test clean lint run audit scan

all: build

build:
	go build $(LDFLAGS) -o $(BINARY_NAME).exe .

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY_NAME).exe

audit: build
	./$(BINARY_NAME).exe audit .github/workflows/ci.yml

scan: build
	./$(BINARY_NAME).exe scan .
