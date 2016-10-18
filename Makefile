build:
	go build -o build/dp-frontend-router

debug: build
	HUMAN_LOG=1 ./build/dp-frontend-router

.PHONY: build debug
