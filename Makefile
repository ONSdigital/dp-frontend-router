SHELL=bash

BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

# --- Production environment
build: generate-build
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"

generate-build:
	cd assets; go-bindata -o redirects.go.new -pkg assets redirects/...
	{ echo "// +build production"; cat assets/redirects.go.new; } > assets/redirects.go
	rm -f assets/redirects.go.new
	rm -f assets/debug.go

# --- Dev environment (please don't commit deletion of assets/redirect.go file!)
debug: generate-debug
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-router

generate-debug:
	cd assets; go-bindata -debug -o debug.go.new -pkg assets redirects/...
	{ echo "// +build debug"; cat assets/debug.go.new; } > assets/debug.go
	rm -f assets/debug.go.new
	rm -f assets/redirects.go

# --- Test
test:
	go test -race -cover -tags 'production' ./...


.PHONY: build debug
