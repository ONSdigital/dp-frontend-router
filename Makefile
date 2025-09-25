SHELL=bash

BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	dis-vulncheck

.PHONY: build
build:
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: debug
debug:
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-router

.PHONY: debug-watch
debug-watch: 
	reflex -d none -c ./reflex

.PHONY: debug-run
debug-run:
	HUMAN_LOG=1 go run -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -race $(LDFLAGS) main.go

.PHONY: test
test:
	go test -race -cover -tags 'production' ./...

.PHONY: test-component
test-component: ## does not run component test. Added as part of nomad pipeline
	exit

