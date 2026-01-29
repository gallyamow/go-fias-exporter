APP_NAME := fias-exporter
APP_MAIN := ./cmd/main.go
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X "main.version=$(VERSION)"

.PHONY: build run clean version

build:
	@echo "Building $(APP_NAME) version $(VERSION)"
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(APP_NAME) $(APP_MAIN)

run:
	@echo "Running $(APP_NAME) version $(VERSION)"
	go run -ldflags '$(LDFLAGS)' .

clean:
	rm -f $(APP_NAME)

test:
	go test  ./...

lint:
	goimports -w .
	golangci-lint run ./...

version:
	@echo $(VERSION)