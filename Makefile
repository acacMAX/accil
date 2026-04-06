.PHONY: help build test clean lint fmt vet install release-windows release-mac release-linux release-all

# Variables
BINARY_NAME=accil
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=dist
INSTALL_DIR=$(USERPROFILE)/.accil/bin

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2 2>/dev/null || echo "")
YELLOW := $(shell tput -Txterm setaf 3 2>/dev/null || echo "")
RED    := $(shell tput -Txterm setaf 1 2>/dev/null || echo "")
NC     := $(shell tput -Txterm sgr0 2>/dev/null || echo "")

help: ## Show this help message
	@echo 'Usage:'
	@echo '  ${GREEN}make${NC} ${YELLOW}<target>${NC}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary for current OS
	@echo "${GREEN}Building ${BINARY_NAME}...${NC}"
	go build -ldflags="-X main.Version=$(VERSION)" -o $(BINARY_NAME).exe .
	@echo "${GREEN}Build complete: ${BINARY_NAME}.exe${NC}"
	@echo "${GREEN}Installing to global...${NC}"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME).exe $(INSTALL_DIR)/
	@echo "${GREEN}Installed to: $(INSTALL_DIR)/$(BINARY_NAME).exe${NC}"

test: ## Run tests
	@echo "${GREEN}Running tests...${NC}"
	go test -v ./...

clean: ## Clean build artifacts
	@echo "${YELLOW}Cleaning...${NC}"
	rm -rf $(BUILD_DIR) $(BINARY_NAME) $(BINARY_NAME).exe
	@echo "${GREEN}Clean complete${NC}"

lint: ## Run linter
	@echo "${GREEN}Running linter...${NC}"
	golangci-lint run ./...

fmt: ## Format code
	@echo "${GREEN}Formatting code...${NC}"
	go fmt ./...

vet: ## Run go vet
	@echo "${GREEN}Running go vet...${NC}"
	go vet ./...

install: ## Install the binary
	@echo "${GREEN}Installing ${BINARY_NAME}...${NC}"
	go install -ldflags="-X main.Version=$(VERSION)" .

release-windows: ## Build for Windows
	@echo "${GREEN}Building for Windows...${NC}"
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "${GREEN}Windows build complete${NC}"

release-mac: ## Build for macOS
	@echo "${GREEN}Building for macOS...${NC}"
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@echo "${GREEN}macOS build complete${NC}"

release-linux: ## Build for Linux
	@echo "${GREEN}Building for Linux...${NC}"
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "${GREEN}Linux build complete${NC}"

release-all: ## Build for all platforms
	@echo "${GREEN}Building for all platforms...${NC}"
	mkdir -p $(BUILD_DIR)
	# Windows
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	# macOS
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	# Linux
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@echo "${GREEN}All builds complete!${NC}"
	@echo ""
	@echo "Binaries are in $(BUILD_DIR)/"
	@ls -lh $(BUILD_DIR)/

dev: ## Run in development mode
	@echo "${GREEN}Running in development mode...${NC}"
	go run .

setup: ## First-time setup
	@echo "${GREEN}Running setup wizard...${NC}"
	go run . --setup

.DEFAULT_GOAL := build
