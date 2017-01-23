build: generate
	go build -tags 'production' -o build/dp-frontend-router

debug: generate
	go build -tags 'debug' -o build/dp-frontend-router
	HUMAN_LOG=1 DEBUG=1 ./build/dp-frontend-router

generate:
	# build the production version
	go generate ./...
	{ echo "// +build production"; cat assets/templates.go; } > assets/templates.go.new
	mv assets/templates.go.new assets/templates.go
	# build the dev version
	cd assets; go-bindata -debug -o debug.go -pkg assets templates/...
	{ echo "// +build debug"; cat assets/debug.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/debug.go

.PHONY: build debug
