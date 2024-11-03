GO_FLAGS	?=
NAME		:= pingreboot
PACKAGE		:= github.com/dabump/$(NAME)
BIN_PREFIX	:= bin

default: help

.PHONY: test
unit-test: ## Run all tests
	@go clean --testcache && go test ./... -v -coverprofile cover.out

generate-mocks: ## Generate mocks
	@go generate ./...

cover: ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build: ## Builds the token bucket
	@go build ${GO_FLAGS} -o ${BIN_PREFIX}/pingreboot ./cmd/pingreboot/main.go

format: ## Format and organise imports
	@go install mvdan.cc/gofumpt@latest
	@gofumpt -l -w .

golangci-lint: ## lint runner
	@go install github.com/golangci/golangci-lint@latest
	@golangci-lint  run --config ./.github/config/golangci.yml -v

gci: ## gci tool to organise imports
	@go install github.com/daixiang0/gci@latest
	@gci write --skip-generated -s default -s standard .

clean: ## Cleans the build binaries
	@rm -rf ${BIN_PREFIX}

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'
