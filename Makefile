build: generate
	go build -tags 'production' -o build/dp-frontend-router

debug: generate
	go build -tags 'debug' -o build/dp-frontend-router
	HUMAN_LOG=1 DEBUG=1 ./build/dp-frontend-router

generate: ${GOPATH}/bin/go-bindata
	# build the production version
	cd assets; ${GOPATH}/bin/go-bindata -o templates.go -pkg assets templates/...
	{ echo "// +build production"; cat assets/templates.go; } > assets/templates.go.new
	mv assets/templates.go.new assets/templates.go
	# build the dev version
	cd assets; ${GOPATH}/bin/go-bindata -debug -o debug.go -pkg assets templates/...
	{ echo "// +build debug"; cat assets/debug.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/debug.go

test:
	go test -tags 'production' ./...

${GOPATH}/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/go-bindata

.PHONY: build debug
