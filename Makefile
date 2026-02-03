.PHONY: build build-web build-all clean

build-web:
	cd web && pnpm install && pnpm run build

build:
	go build -o bin/goyavision ./cmd/server

build-all: build-web build

clean:
	rm -rf web/dist web/node_modules bin
