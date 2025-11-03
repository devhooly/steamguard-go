.PHONY: build clean install test run help

# Variables
BINARY_NAME=steamguard
MAIN_PATH=.
INSTALL_PATH=/usr/local/bin

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

help: ## Show help
	@echo "$(GREEN)Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

build: ## Build binary
	@echo "$(GREEN)Building $(BINARY_NAME)...$(NC)"
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Done!$(NC)"

build-all: ## Build for all platforms
	@echo "$(GREEN)Building for all platforms...$(NC)"
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "$(GREEN)✓ Done! Files in dist/$(NC)"

clean: ## Clean built files
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -f $(BINARY_NAME)
	rm -rf dist/
	go clean
	@echo "$(GREEN)✓ Cleaned$(NC)"

install: build ## Install to system
	@echo "$(GREEN)Installing $(BINARY_NAME) to $(INSTALL_PATH)...$(NC)"
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/
	@echo "$(GREEN)✓ Installed!$(NC)"

uninstall: ## Uninstall from system
	@echo "$(YELLOW)Removing $(BINARY_NAME)...$(NC)"
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)✓ Removed$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test -v ./...

run: build ## Build and run
	@echo "$(GREEN)Running $(BINARY_NAME)...$(NC)"
	./$(BINARY_NAME)

deps: ## Install dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	go mod download
	go mod verify
	@echo "$(GREEN)✓ Done!$(NC)"

tidy: ## Tidy go.mod
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	go mod tidy
	@echo "$(GREEN)✓ Done!$(NC)"

fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✓ Done!$(NC)"

lint: ## Check code with linter
	@echo "$(GREEN)Linting code...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(YELLOW)Install golangci-lint: https://golangci-lint.run/usage/install/$(NC)" && exit 1)
	golangci-lint run
	@echo "$(GREEN)✓ Linting complete$(NC)"

dev: ## Development mode with auto-reload
	@which air > /dev/null || (echo "$(YELLOW)Install air: go install github.com/cosmtrek/air@latest$(NC)" && exit 1)
	air
