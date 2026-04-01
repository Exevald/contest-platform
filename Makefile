SHELL := /bin/zsh

APP_NAME := contestplatform
FRONTEND_DIR := frontend
BUILD_DIR := build/bin
GO_CACHE_DIR := $(CURDIR)/.cache/go-build
GO_MOD_CACHE_DIR := $(CURDIR)/.cache/gomod
GO_ENV := GOCACHE=$(GO_CACHE_DIR) GOMODCACHE=$(GO_MOD_CACHE_DIR)
WAILS ?= wails
COMPILERS_DIR ?=

.PHONY: help setup deps frontend-deps frontend-build frontend-clean test dev build build-clean build-macos build-linux build-windows cross-build release-check clean

help:
	@echo "ContestPlatform build targets"
	@echo ""
	@echo "Setup"
	@echo "  make setup            Install frontend deps and prepare local Go caches"
	@echo "  make deps             Install frontend deps"
	@echo ""
	@echo "Development"
	@echo "  make dev              Run Wails in dev mode"
	@echo "  make test             Run Go tests"
	@echo "  make frontend-build   Build frontend with Vite"
	@echo ""
	@echo "Production build"
	@echo "  make build            Build the app with 'wails build'"
	@echo "  make build-macos      Build macOS bundle with Wails"
	@echo "  make build-linux      Build Linux package with Wails"
	@echo "  make build-windows    Build Windows package with Wails"
	@echo ""
	@echo "Cross-check binaries"
	@echo "  make cross-build      Build tagged production binaries for darwin/linux/windows"
	@echo "  make release-check    Run tests, frontend build and cross-build"
	@echo ""
	@echo "Cleanup"
	@echo "  make clean            Remove local caches and generated frontend dist"
	@echo "  make build-clean      Remove build/bin contents"
	@echo ""
	@echo "Notes"
	@echo "  1. Run the packaged Wails output from build/bin, for example contestplatform.app on macOS."
	@echo "  2. Do not run raw binaries created by plain 'go build' without the 'production' tag."
	@echo "  3. To use bundled compilers, set COMPILERS_DIR=/absolute/path/to/compilers before build/run."

setup: deps
	@mkdir -p "$(GO_CACHE_DIR)" "$(GO_MOD_CACHE_DIR)" "$(BUILD_DIR)"

deps: frontend-deps

frontend-deps:
	cd "$(FRONTEND_DIR)" && npm install

frontend-build:
	cd "$(FRONTEND_DIR)" && npm run build

frontend-clean:
	rm -rf "$(FRONTEND_DIR)/dist"

test: setup
	$(GO_ENV) go test ./...

dev: setup
	COMPILERS_DIR="$(COMPILERS_DIR)" $(WAILS) dev

build: setup
	COMPILERS_DIR="$(COMPILERS_DIR)" $(WAILS) build

build-macos: setup
	COMPILERS_DIR="$(COMPILERS_DIR)" $(WAILS) build -platform darwin/universal

build-linux: setup
	COMPILERS_DIR="$(COMPILERS_DIR)" $(WAILS) build -platform linux/amd64

build-windows: setup
	COMPILERS_DIR="$(COMPILERS_DIR)" $(WAILS) build -platform windows/amd64

cross-build: setup frontend-build build-clean
	$(GO_ENV) GOOS=darwin GOARCH=arm64 go build -tags production -o "$(BUILD_DIR)/$(APP_NAME)-darwin-arm64" .
	$(GO_ENV) GOOS=linux GOARCH=amd64 go build -tags production -o "$(BUILD_DIR)/$(APP_NAME)-linux-amd64" .
	$(GO_ENV) GOOS=windows GOARCH=amd64 go build -tags production -o "$(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe" .

release-check: test frontend-build cross-build

build-clean:
	rm -rf "$(BUILD_DIR)"
	mkdir -p "$(BUILD_DIR)"

clean: build-clean frontend-clean
	rm -rf "$(CURDIR)/.cache"
