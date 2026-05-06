VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build install test clean

build:
	go build $(LDFLAGS) -o bin/chameleon .

install:
	go install $(LDFLAGS) .

test:
	go test ./... -v

clean:
	rm -rf bin/
