BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

build: generate
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"

debug: generate
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router -ldflags="-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)' -X 'main.Version=$(VERSION)'"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-router

generate: ${GOPATH}/bin/go-bindata
	# build the production version
	cd assets; ${GOPATH}/bin/go-bindata -o redirects.go -pkg assets redirects/...
	{ echo "// +build production"; cat assets/redirects.go; } > assets/redirects.go.new
	mv assets/redirects.go.new assets/redirects.go
	# build the dev version
	cd assets; ${GOPATH}/bin/go-bindata -debug -o debug.go -pkg assets redirects/...
	{ echo "// +build debug"; cat assets/debug.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/debug.go

test:
	go test -tags 'production' ./...

${GOPATH}/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/go-bindata

.PHONY: build debug
