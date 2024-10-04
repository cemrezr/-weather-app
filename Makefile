APP_NAME = weather-app
CMD_DIR = cmd

GO_FILES = $(shell find . -name '*.go' -not -path "./vendor/*")
GO_BUILD_DIR = ./bin

all: build

deps:
	@echo "==> Downloading dependencies..."
	go mod tidy

build: deps
	@echo "==> Building the application..."
	mkdir -p $(GO_BUILD_DIR)
	go build -o $(GO_BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

run: build
	@echo "==> Running the application..."
	$(GO_BUILD_DIR)/$(APP_NAME)

clean:
	@echo "==> Cleaning build cache and output files..."
	rm -rf $(GO_BUILD_DIR)

test:
	@echo "==> Running tests..."
	go test ./...

lint:
	@echo "==> Running golangci-lint..."
	golangci-lint run


.PHONY: all deps build run clean test lint help