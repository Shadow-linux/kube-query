NAME := kube-query
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.Version=$(VERSION)' \
           -X 'main.Revision=$(REVISION)'
GO ?= GO111MODULE=on go
.DEFAULT_GOAL := help

.PHONY: darwin-amd64
darwin-amd64: ## Build a darwin-amd64 pkg.
	@# darwin-amd64
	rm -f bin/$(NAME).darwin-amd64 2> /dev/null; \
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(NAME).darwin-amd64 main.go; \
	zip pkg/$(NAME)_$(VERSION)_darwin_amd64.zip bin/$(NAME).darwin-amd64;

.PHONY: darwin-arm64
darwin-arm64: ## Build a darwin-arm64 pkg.
	rm -f bin/$(NAME).darwin-arm64 2> /dev/null; \
    GOOS=darwin GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(NAME).darwin-arm64 main.go; \
    zip pkg/$(NAME)_$(VERSION)_darwin_arm64.zip bin/$(NAME).darwin-arm64;

.PHONY: linux-amd64
linux-amd64: ## Build a linux-amd64 pkg.
	# linux
	rm -f bin/$(NAME).linux-amd64 2> /dev/null; \
    GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(NAME).linux-amd64 main.go; \
    zip pkg/$(NAME)_$(VERSION)_linux_amd64.zip bin/$(NAME).linux-amd64;


.PHONY: linux-arm64
linux-arm64: ## Build a linux-arm64 pkg.
	# linux
	rm -f bin/$(NAME).linux-arm64 2> /dev/null; \
    GOOS=linux GOARCH=arm64 $(GO) build -ldflags "$(LDFLAGS)" -o bin/$(NAME).linux-arm64 main.go; \
    zip pkg/$(NAME)_$(VERSION)_linux_arm64.zip bin/$(NAME).linnux-arm64;

.PHONY: chmod
chmod:
	chmod a+x bin/*

.PHONY: build
build: main.go  ## Build a binary.
	$(GO) build -ldflags "$(LDFLAGS)" -o bin/$(NAME) main.go

.PHONY: help
help: ## Show help text
	@echo "Commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'