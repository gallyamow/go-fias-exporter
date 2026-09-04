APP_NAME := fias-exporter
APP_MAIN := ./cmd/main.go
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X "main.version=$(VERSION)"

.PHONY: build run clean test lint version

build:
	@echo "Building $(APP_NAME) version $(VERSION)"
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(APP_NAME) $(APP_MAIN)

build-linux-amd64:
	@echo "Building Linux amd64"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags '$(LDFLAGS)' -o $(APP_NAME)-linux-amd64 $(APP_MAIN)

build-darwin-amd64:
	@echo "Building macOS amd64"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 \
		go build -ldflags '$(LDFLAGS)' -o $(APP_NAME)-darwin-amd64 $(APP_MAIN)

build-darwin-arm64:
	@echo "Building macOS arm64"
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 \
		go build -ldflags '$(LDFLAGS)' -o $(APP_NAME)-darwin-arm64 $(APP_MAIN)

run:
	@echo "Running $(APP_NAME) version $(VERSION)"
	go run -ldflags '$(LDFLAGS)' .

clean:
	rm -f $(APP_NAME) $(APP_NAME)-linux-amd64 $(APP_NAME)-darwin-amd64 $(APP_NAME)-darwin-arm64

test:
	go test ./...

lint:
	goimports -w .
	golangci-lint run ./...

version:
	@echo $(VERSION)