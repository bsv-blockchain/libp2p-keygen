.PHONY: build clean install build-all test

# Binary name
BINARY=libp2p-keygen
VERSION=1.0.0

# Build output directory
BUILD_DIR=bin

# Go build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

# Default target - build for current platform
build:
	@echo "Building $(BINARY) for current platform..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY)"

# Install to $GOPATH/bin or $HOME/go/bin
install:
	@echo "Installing $(BINARY)..."
	go install $(LDFLAGS) .
	@echo "Installation complete"

# Build for all platforms
build-all: clean
	@echo "Building $(BINARY) for all platforms..."
	@mkdir -p $(BUILD_DIR)

	# macOS (Intel)
	@echo "Building for macOS (amd64)..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64 .

	# macOS (Apple Silicon)
	@echo "Building for macOS (arm64)..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64 .

	# Linux (amd64)
	@echo "Building for Linux (amd64)..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64 .

	# Linux (arm64)
	@echo "Building for Linux (arm64)..."
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64 .

	# Windows (amd64)
	@echo "Building for Windows (amd64)..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .

	@echo "Build complete! Binaries are in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Quick test run
run:
	go run . $(ARGS)

# Example runs
examples:
	@echo "Example: Simple generation"
	go run .
	@echo "\nExample: Vanity generation (will take a moment)"
	go run . --vanity test

# Show help
help:
	@echo "Available targets:"
	@echo "  make build       - Build for current platform"
	@echo "  make build-all   - Build for all platforms (darwin, linux, windows)"
	@echo "  make install     - Install to GOPATH/bin"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Remove build artifacts"
	@echo "  make run         - Run directly (use ARGS='--vanity foo' for arguments)"
	@echo "  make examples    - Run example commands"
	@echo "  make help        - Show this help message"
