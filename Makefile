BINARY_NAME=opsdev
VERSION=1.0.0
BUILD_DIR=build
MAIN_PATH=cmd/opsdev/main.go

.PHONY: all build clean install test

all: clean build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo mv $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

release: build-all
	@echo "Creating release package..."
	@cd $(BUILD_DIR) && \
	tar czf $(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	tar czf $(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar czf $(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64