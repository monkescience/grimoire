.PHONY: build test lint fmt clean install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o grimoire ./cmd/grimoire

test:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .
	goimports -w .

clean:
	rm -f grimoire

install:
	go install $(LDFLAGS) ./cmd/grimoire
