SHELL=bash

BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

build:
	cd assets; go-bindata -o redirects.go -pkg assets redirects/...
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"

debug:
	cd assets; go-bindata -debug -o redirects.go -pkg assets redirects/...
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-router

test:
	go test -race -cover -tags 'production' ./...


.PHONY: build debug
