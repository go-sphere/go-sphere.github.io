GO ?= go
GOLANGCI_LINT ?= golangci-lint
HUGO ?= hugo

.DEFAULT_GOAL := check

.PHONY: deps-update tidy fmt test lint check

# The module's only dependency is the Hugo theme (github.com/imfing/hextra),
# which ships no Go packages—so `go mod tidy` would drop it. Use Hugo's module
# commands, which read module.imports from hugo.yaml.
deps-update:
	$(HUGO) mod get -u
	$(HUGO) mod tidy

tidy:
	$(HUGO) mod tidy

fmt:
	$(GO) fmt ./...
	$(GOLANGCI_LINT) fmt --no-config --enable gofmt --enable goimports

test:
	$(GO) test ./...

lint:
	$(GOLANGCI_LINT) fmt --no-config --enable gofmt --enable goimports --diff
	$(GO) vet ./...
	$(GOLANGCI_LINT) run --no-config

check:
	$(HUGO) mod tidy
	@git diff --exit-code -- go.mod go.sum || { echo "go.mod/go.sum are not tidy; run 'make tidy'"; exit 1; }
	$(MAKE) lint
	$(MAKE) test
