# ==============================================================================
# 🎬 Vidonex - Declarative Video Composition & FFmpeg Filtergraph Engine
# Modern, High-Performance Build System & Developer Automation
# ==============================================================================

SHELL := /usr/bin/env bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := help

# ─── Variables & Configuration ────────────────────────────────────────────────
APP_NAME            := vidonex
CLI_BINARY          := bin/vidonex
DESKTOP_BINARY      := bin/vidonex-desktop
VERSION             ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "1.0.0")
COMMIT              ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME          ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_LDFLAGS          := -s -w -X 'main.EngineVersion=$(VERSION)' -X 'main.GitCommit=$(COMMIT)' -X 'main.BuildTime=$(BUILD_TIME)'

# Colors for modern, polished terminal output
CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
BLUE   := \033[34m
BOLD   := \033[1m
RESET  := \033[0m

# ─── Help Target ──────────────────────────────────────────────────────────────
.PHONY: help
help: ## Display this colorful and organized command reference
	@printf "\n$(BOLD)$(CYAN)🎬 Vidonex Build & Development Commands $(RESET) (v$(VERSION))\n\n"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / { \
		category = "General"; \
		if ($$1 ~ /^(desktop|dev-desktop)/) category = "Desktop Workstation"; \
		else if ($$1 ~ /^(build|cli)/) category = "Building & Binaries"; \
		else if ($$1 ~ /^(test|lint|cover)/) category = "Testing & Quality"; \
		else if ($$1 ~ /^(ui|generate)/) category = "Web Studio Frontend"; \
		else if ($$1 ~ /^(run|serve)/) category = "Runtime & Studio"; \
		else if ($$1 ~ /^(clean|deps|install)/) category = "Maintenance & Setup"; \
		printf "  $(BOLD)$(GREEN)%-20s$(RESET) %s $(YELLOW)[%s]$(RESET)\n", $$1, $$2, category \
	}' $(MAKEFILE_LIST) | sort -k3
	@printf "\n$(BOLD)Quick Start:$(RESET)\n"
	@printf "  $(CYAN)make run-desktop$(RESET)      Launch native desktop workstation\n"
	@printf "  $(CYAN)make dev-desktop$(RESET)      Start live desktop development with Hot-Reload\n"
	@printf "  $(CYAN)make serve$(RESET)            Start web studio server\n"
	@printf "  $(CYAN)make test$(RESET)             Run all unit and E2E test suites\n\n"

# ─── Desktop Workstation (Wails v2 + Nuxt 4) ──────────────────────────────────
.PHONY: desktop
desktop: ui-build install-desktop-icon ## Compile the native single-binary desktop workstation
	@printf "$(CYAN)🔨 Building native desktop binary ($(DESKTOP_BINARY))...$(RESET)\n"
	@mkdir -p bin
	@go build -trimpath -tags "webkit2_41,desktop,production" -ldflags="$(GO_LDFLAGS)" -o $(DESKTOP_BINARY) .
	@printf "$(GREEN)✅ Desktop binary compiled successfully: $(BOLD)%s$(RESET) (%s)\n" \
		"$(DESKTOP_BINARY)" "$$(du -h $(DESKTOP_BINARY) | cut -f1)"

.PHONY: run-desktop
run-desktop: desktop ## Build and immediately run the native desktop application
	@printf "$(CYAN)🚀 Launching Vidonex Desktop Workstation...$(RESET)\n"
	@./$(DESKTOP_BINARY)

.PHONY: install-desktop-icon
install-desktop-icon: ## Register desktop icon and launcher for Linux desktop integration
	@if [ "$$(uname -s)" = "Linux" ]; then \
		mkdir -p $(HOME)/.local/share/icons/hicolor/512x512/apps $(HOME)/.local/share/pixmaps $(HOME)/.local/share/applications $(HOME)/.local/bin; \
		cp -f build/appicon.png $(HOME)/.local/share/icons/hicolor/512x512/apps/vidonex.png 2>/dev/null || true; \
		cp -f build/appicon.png $(HOME)/.local/share/pixmaps/vidonex.png 2>/dev/null || true; \
		cp -f build/linux/vidonex.desktop $(HOME)/.local/share/applications/vidonex.desktop 2>/dev/null || true; \
		ln -sf $(CURDIR)/bin/vidonex-desktop $(HOME)/.local/bin/vidonex-desktop 2>/dev/null || true; \
		which gtk-update-icon-cache >/dev/null 2>&1 && gtk-update-icon-cache -f -t $(HOME)/.local/share/icons/hicolor 2>/dev/null || true; \
		which update-desktop-database >/dev/null 2>&1 && update-desktop-database $(HOME)/.local/share/applications 2>/dev/null || true; \
	fi

.PHONY: dev-desktop
dev-desktop: ## Start desktop workstation in live development mode with Hot-Reload (requires wails)
	@printf "$(CYAN)⚡ Starting Wails live development mode (Hot-Reload)...$(RESET)\n"
	@if command -v wails >/dev/null 2>&1; then \
		wails dev -s -tags webkit2_41; \
	elif [ -f $(HOME)/go/bin/wails ]; then \
		$(HOME)/go/bin/wails dev -s -tags webkit2_41; \
	else \
		printf "$(RED)Wails CLI not found. Run 'make install-wails' first.$(RESET)\n"; exit 1; \
	fi

.PHONY: install-wails
install-wails: ## Install or update Wails v2 CLI binary
	@printf "$(CYAN)📦 Installing Wails v2 CLI to GOPATH...$(RESET)\n"
	@go install github.com/wailsapp/wails/v2/cmd/wails@v2.15.0
	@printf "$(GREEN)✅ Wails installed at: $(HOME)/go/bin/wails$(RESET)\n"

# ─── Web Studio Frontend (Nuxt 4 + Nuxt UI) ───────────────────────────────────
.PHONY: ui-deps
ui-deps: ## Install frontend node dependencies
	@printf "$(CYAN)📦 Installing UI frontend packages...$(RESET)\n"
	@npm --prefix ui install

.PHONY: ui-dev
ui-dev: ## Start Nuxt 4 frontend standalone dev server (http://localhost:3000)
	@printf "$(CYAN)⚡ Launching Nuxt 4 Web Studio dev server...$(RESET)\n"
	@npm --prefix ui run dev

.PHONY: ui-build
ui-build: ## Build and generate static UI distribution for embedding
	@printf "$(CYAN)📦 Generating static UI distribution (public/dist)...$(RESET)\n"
	@rm -rf ui/public/dist
	@npm --prefix ui run generate

# ─── Building & Binaries ──────────────────────────────────────────────────────
.PHONY: build
build: cli desktop ## Build both the standalone CLI and Desktop application

.PHONY: cli
cli: ## Compile the Vidonex standalone CLI binary (bin/vidonex)
	@printf "$(CYAN)🔨 Compiling standalone CLI binary ($(CLI_BINARY))...$(RESET)\n"
	@mkdir -p bin
	@go build -trimpath -ldflags="$(GO_LDFLAGS)" -o $(CLI_BINARY) ./cmd/vidonex
	@printf "$(GREEN)✅ CLI binary compiled: $(BOLD)%s$(RESET)\n" "$(CLI_BINARY)"

.PHONY: wasm
wasm: ## Compile the Go engine to WebAssembly (ui/public/vidonex.wasm)
	@printf "$(CYAN)🌐 Compiling Vidonex engine to WebAssembly...$(RESET)\n"
	@mkdir -p ui/public
	@GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o ui/public/vidonex.wasm ./cmd/wasm
	@cp "$$(go env GOROOT)/lib/wasm/wasm_exec.js" ui/public/wasm_exec.js 2>/dev/null || cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" ui/public/wasm_exec.js 2>/dev/null || true
	@printf "$(GREEN)✅ WebAssembly binary compiled: $(BOLD)ui/public/vidonex.wasm$(RESET) (%s)\n" "$$(du -h ui/public/vidonex.wasm | cut -f1)"

# ─── Runtime & Execution ──────────────────────────────────────────────────────
.PHONY: serve
serve: cli ## Start the embedded HTTP/WebSocket server and serve Web Studio
	@printf "$(CYAN)📡 Starting Vidonex Web Studio Server...$(RESET)\n"
	@./$(CLI_BINARY) serve --port 8080

.PHONY: probe
probe: cli ## Probe a media file metadata (usage: make probe FILE=clip.mp4)
	@if [ -z "$(FILE)" ]; then \
		printf "$(RED)Usage: make probe FILE=<path_to_media>$(RESET)\n"; exit 1; \
	fi
	@./$(CLI_BINARY) probe "$(FILE)"

.PHONY: validate
validate: cli ## Validate a declarative YAML spec (usage: make validate FILE=project.yaml)
	@if [ -z "$(FILE)" ]; then \
		printf "$(RED)Usage: make validate FILE=<path_to_yaml>$(RESET)\n"; exit 1; \
	fi
	@./$(CLI_BINARY) validate "$(FILE)"

.PHONY: graph
graph: cli ## Export filtergraph Mermaid diagram (usage: make graph FILE=project.yaml)
	@if [ -z "$(FILE)" ]; then \
		printf "$(RED)Usage: make graph FILE=<path_to_yaml>$(RESET)\n"; exit 1; \
	fi
	@./$(CLI_BINARY) graph "$(FILE)"

# ─── Testing & Quality Assurance ──────────────────────────────────────────────
.PHONY: test
test: ## Run all unit and feature tests with race detector
	@printf "$(CYAN)🧪 Running test suite with Race Detector...$(RESET)\n"
	@go test -race -v ./...

.PHONY: test-e2e
test-e2e: ## Run real-FFmpeg end-to-end integration tests
	@printf "$(CYAN)🎥 Running Real-FFmpeg E2E tests...$(RESET)\n"
	@go test -race -v -timeout=300s ./tests/e2e/...

.PHONY: lint
lint: ## Run golangci-lint on all Go packages
	@printf "$(CYAN)🔍 Running golangci-lint...$(RESET)\n"
	@golangci-lint run ./...
	@printf "$(GREEN)✅ Linter passed cleanly!$(RESET)\n"

.PHONY: coverage
coverage: ## Generate test coverage profile and display summary
	@printf "$(CYAN)📊 Calculating test coverage...$(RESET)\n"
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -func=coverage/coverage.out | tail -n 1
	@printf "$(GREEN)To view in browser: go tool cover -html=coverage/coverage.out$(RESET)\n"

# ─── Maintenance & Cleanup ────────────────────────────────────────────────────
.PHONY: deps
deps: ## Download and tidy Go modules
	@printf "$(CYAN)📦 Tidying Go dependencies...$(RESET)\n"
	@go mod tidy
	@go mod verify

.PHONY: clean
clean: ## Clean build artifacts, temp files, and caches
	@printf "$(YELLOW)🧹 Cleaning build artifacts and cache directories...$(RESET)\n"
	@rm -rf bin/ coverage/ build/bin/
	@rm -rf ui/.output ui/.nuxt ui/dist
	@go clean -cache -testcache
	@printf "$(GREEN)✨ Clean complete!$(RESET)\n"

# ─── Documentation (VitePress) ────────────────────────────────────────────────
.PHONY: docs-install
docs-install: ## Install documentation site dependencies
	@printf "$(CYAN)📦 Installing documentation dependencies...$(RESET)\n"
	@cd docs && pnpm install --dangerously-allow-all-builds

.PHONY: docs-dev
docs-dev: ## Start live-reloading documentation dev server
	@printf "$(CYAN)📖 Starting documentation dev server on http://localhost:5173/vidonex/...$(RESET)\n"
	@cd docs && pnpm run dev

.PHONY: docs-build
docs-build: ## Build production static documentation site for GitHub Pages
	@printf "$(CYAN)🏗️ Building static documentation site...$(RESET)\n"
	@cd docs && pnpm run build
	@printf "$(GREEN)✅ Documentation built in docs/.vitepress/dist!$(RESET)\n"
