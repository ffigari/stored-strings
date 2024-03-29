all: build test

build: build-web-api

build-web-api:
	go build -o binwebapi cmd/webapi/main.go

test:
	go test ./...

clean:
	rm -f binwebapi

run: build
	./binwebapi
