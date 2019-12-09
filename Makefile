default: build

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

build: always
	GO111MODULE=on
	go build -v -ldflags="-s -w" -o .target/dice cmd/dice/main.go

.PHONY: clean
clean:
    rm -rf .target

.PHONY: test
test:
	go test -v ./...
