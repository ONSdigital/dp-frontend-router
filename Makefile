SHELL=bash

BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: assets
assets:
	# cd assets; go run github.com/jteeuwen/go-bindata/go-bindata -o redirects.go -pkg assets redirects/...

.PHONY: assets-debug
assets-debug:
	# cd assets; go run github.com/jteeuwen/go-bindata/go-bindata -debug -o redirects.go -pkg assets redirects/...

.PHONY: clean-assets
clean-assets:
	rm assets/redirects.go

.PHONY: build
build: assets
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"

.PHONY: lint
lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	golangci-lint run ./...

.PHONY: debug
debug: assets-debug
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-router

.PHONY: debug-run
debug-run: assets-debug
	HUMAN_LOG=1 go run -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)" -race $(LDFLAGS) main.go
	clean-assets

.PHONY: test
test: assets
	go test -race -cover -tags 'production' ./...

.PHONY: test-component
test-component: ## does not run component test. Added as part of nomad pipeline
	exit

