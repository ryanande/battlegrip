.PHONY: default help

## this needs to be refactored and cleaned up, way too much going on here

default: help
help: ## help: display make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m make %-20s -> %s\n\033[0m", $$1, $$2}'
	
# make: app info
APP_NAME    := battlegrip
APP_WORKDIR := $(shell pwd)
CLIENT_ID   ?= bogus.account
APP_LOG_FMT := `/bin/date "+%Y-%m-%d %H:%M:%S %z [$(APP_NAME)]"`
	
# make: go info
GO_COV_DIR  := $(APP_WORKDIR)/coverage
GO_UNIT_DIR := $(GO_COV_DIR)/unit
GO_INT_DIR := $(GO_COV_DIR)/integration
GO_PACKAGES := $(shell go list -f '{{.Dir}}' ./...)
	
# make: go test info
GO_UNIT_JUNIT     := $(GO_UNIT_DIR)/junit.xml
GO_UNIT_WEBPAGE   := $(GO_UNIT_DIR)/index.html
GO_UNIT_REPORT    := $(GO_UNIT_DIR)/report.out
GO_UNIT_COVERAGE  := $(GO_UNIT_DIR)/coverage.out
GO_UNIT_COBERTURA := $(GO_UNIT_DIR)/cobertura.xml
	
GO_INT_JUNIT     := $(GO_INT_DIR)/junit.xml
GO_INT_WEBPAGE   := $(GO_INT_DIR)/index.html
GO_INT_REPORT    := $(GO_INT_DIR)/report.out
GO_INT_COVERAGE  := $(GO_INT_DIR)/coverage.out
GO_INT_COBERTURA := $(GO_INT_DIR)/cobertura.xml

	
# --------------------------------------------------
# Build Targets
# --------------------------------------------------
.PHONY: vendor
vendor:
	@go mod vendor 
# --------------------------------------------------
# Build Targets
# --------------------------------------------------
.PHONY: build
build:
	@go build . 
# --------------------------------------------------
# Test Targets
# --------------------------------------------------
.PHONY: test-clean
test-clean: ## test: clean workspace
	@rm -rf $(GO_COV_DIR)


.PHONY: lint-go
lint-go: ## lint: lints golang and tries to automatically fix things
	@echo $(APP_LOG_FMT) "linting and trying to fix things automatically"
	@golangci-lint run --fix -v --timeout 3m
	
.PHONY: _test_variables
_test_variables:
	@echo $(APP_LOG_FMT) "installing test dependencies"
	@go install golang.org/x/lint/golint@latest
	@go install github.com/jstemmer/go-junit-report@latest
	@go install github.com/t-yuki/gocover-cobertura@latest
	@echo $(APP_LOG_FMT) "verifying required variables are set"
# ifeq ($(CLIENT_SECRET),)
# 	@echo $(APP_LOG_FMT) "error: CLIENT_SECRET is required"
# 	@exit 1
# endif
	
.PHONY: test-unit
test-unit: _test_variables ## test: run unit tests
	@echo $(APP_LOG_FMT) "running unit test suite"
	@mkdir -p $(GO_UNIT_DIR)
	@go test -v \
	-covermode=atomic \
	-coverprofile=$(GO_UNIT_COVERAGE) \
	$(GO_PACKAGES) \
	2>&1 > $(GO_UNIT_REPORT) || cat $(GO_UNIT_REPORT)
	@cat $(GO_UNIT_REPORT)
	@go tool cover -func=$(GO_UNIT_COVERAGE)
	@go tool cover -html=$(GO_UNIT_COVERAGE) -o $(GO_UNIT_WEBPAGE)
	@cat $(GO_UNIT_REPORT) | go-junit-report > $(GO_UNIT_JUNIT)
	@gocover-cobertura < $(GO_UNIT_COVERAGE) > $(GO_UNIT_COBERTURA)
	
.PHONY: test-integration
test-integration: _test_variables ## test: run integration tests
	@echo $(APP_LOG_FMT) "running integration test suite"
	@mkdir -p $(GO_INT_DIR)
	@CLIENT_ID=$(CLIENT_ID) go test -v \
	-tags=integration \
	-covermode=atomic \
	-coverprofile=$(GO_INT_COVERAGE) \
	$(GO_PACKAGES) \
	2>&1 > $(GO_INT_REPORT) || cat $(GO_INT_REPORT)
	@cat $(GO_INT_REPORT)
	@go tool cover -func=$(GO_INT_COVERAGE)
	@go tool cover -html=$(GO_INT_COVERAGE) -o $(GO_INT_WEBPAGE)
	@cat $(GO_INT_REPORT) | go-junit-report > $(GO_INT_JUNIT)
	@gocover-cobertura < $(GO_INT_COVERAGE) > $(GO_INT_COBERTURA)
	
.PHONY: test-acceptance
test-acceptance: lint-go test-unit test-integration ## test: run acceptance tests
	