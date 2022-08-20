VERSION=$(shell git rev-parse --short HEAD)

.PHONY: build
build:
	mkdir -p ./build

	GOOS=darwin GOARCH=amd64 \
	go build \
	-ldflags "-w -s" \
	-o ./build/bridge-darwin-amd64

	GOOS=linux GOARCH=amd64 \
	go build \
	-ldflags "-w -s" \
	-o ./build/bridge-linux-amd64

	GOOS=linux GOARCH=386 \
	go build \
	-ldflags "-w -s" \
	-o ./build/bridge-linux-x86


	GOOS=windows GOARCH=amd64 \
	go build \
	-ldflags "-w -s" \
	-o ./build/bridge-windows-amd64.exe

