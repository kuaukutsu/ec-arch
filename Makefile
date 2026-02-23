GO_VERSION = "1.26.0"
GOLANGCI_LINT_VERSION = "v2.9.0"
GOFUMPT_VERSION = "v0.9.2"
GOBIN ?= "$$(pwd)/bin"

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
.DEFAULT_GOAL := help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	
goset:
	asdf set golang $(GO_VERSION)

lint: _golangci_lint_install ## golangci-lint
	$(GOBIN)/golangci-lint run ./... --config=.golangci.yml --fix
	
fmt: _gofumpt_install ## gofumpt
	$(GOBIN)/gofumpt -extra -l -w .
	
run: fmt ## Run gRPC server
	HTTP_SERVER_PASSWORD=123456 go run cmd/app/main.go --config=./config/local.yaml

tests: ## Run Tests
	go test ./internal/...
	
_golangci_lint_install:
	[ -f $(GOBIN)/golangci-lint ] || GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

_gofumpt_install:
	[ -f $(GOBIN)/gofumpt ] || GOBIN=$(GOBIN) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
