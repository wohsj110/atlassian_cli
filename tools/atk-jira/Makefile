.PHONY: build test clean install lint tidy check

BINARY_NAME=atk-jira
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags "-X github.com/wohsj110/atlassian_cli/shared/version.Version=$(VERSION) \
	-X github.com/wohsj110/atlassian_cli/shared/version.Commit=$(COMMIT) \
	-X github.com/wohsj110/atlassian_cli/shared/version.BuildDate=$(BUILD_DATE)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/atk-jira

test:
	go test -race -v ./...

test-cover:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

install: build
	cp bin/$(BINARY_NAME) /usr/local/bin/

lint:
	golangci-lint run

tidy:
	go mod tidy
	git diff --exit-code go.mod go.sum

# CI gate: everything that must pass before merge
check: tidy lint test build

deps:
	go mod download

all: check
