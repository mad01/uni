# Build variables
BINARY_NAME=uni
GIT_SHORT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "true" || echo "false")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Linker flags
LDFLAGS = -s -w \
	-X github.com/mad01/uni/cmd.gitHash=$(GIT_SHORT_COMMIT) \
	-X github.com/mad01/uni/cmd.dirty=$(GIT_DIRTY) \
	-X github.com/mad01/uni/cmd.date=$(BUILD_DATE)

.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)

.PHONY: build-release
build-release:
	@echo "Building release version..."
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)

.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install -ldflags "$(LDFLAGS)"

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out
	rm -f coverage.html
	rm -rf dist/

.PHONY: version
version:
	@echo "Git Short: $(GIT_SHORT_COMMIT)"
	@echo "Git Dirty: $(GIT_DIRTY)"
	@echo "Build Date: $(BUILD_DATE)"

.PHONY: run
run: build
	./$(BINARY_NAME)

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

.PHONY: lint
lint: fmt vet

.PHONY: dev
dev: lint test build

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary with git info"
	@echo "  build-release  - Build release binary (CGO disabled)"
	@echo "  install        - Install the binary to GOPATH/bin"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  version        - Show version information"
	@echo "  run            - Build and run the binary"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run fmt and vet"
	@echo "  dev            - Run lint, test, and build"
	@echo "  help           - Show this help" 