.PHONY: help compile build-local build-ssh test fmt clean all

help:
	@echo "Available targets:"
	@echo "compile     - Build all binaries"
	@echo "build-local - Build local TUI app"
	@echo "build-ssh   - Build SSH server app"
	@echo "test        - Run tests"
	@echo "fmt         - Format all go files"
	@echo "clear       - Clear bin/"
	@echo "all         - Clear bin/, run tests, compile"


compile: build-local build-ssh

build-local:
	go build -o bin/app-local ./cmd/app-local/

build-ssh:
	go build -o bin/app-ssh ./cmd/app-ssh/

test:
	go test -race -coverprofile=coverage.out ./...

fmt:
	go fmt ./...

clean:
	rm -rf bin/
	rm -f coverage.out

all: fmt test compile
