BINPATH ?= build

build: generate
	go build -tags 'production' -o $(BINPATH)/dp-frontend-router

debug: generate
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-router
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
