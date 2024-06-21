PROJ = gambol
ORG_PATH = github.com/nuccitheboss
REPO_PATH = $(ORG_PATH)/$(PROJ)

export GOBIN=$(CURDIR)/bin

##@ Build

build: bin/gambol ## Build executable binary

bin/gambol:
	@mkdir -p bin/
	@go install -v $(REPO_PATH)/cmd/gambol

.PHONY: snap
snap: ## Build snap package
	@cp -ra build/snap snap
	@snapcraft -v pack
	@rm -rf snap

##@ Lint

.PHONY: fmt
fmt: ## Format source code
	@go fmt -x ./...
	@golangci-lint run --fix

.PHONY: lint
lint: ## Lint source code
	@golangci-lint run

##@ Test

define E2E_TESTS
	e2e-simple
	e2e-hostpath
	e2e-advanced
endef

.PHONY: e2e
e2e: $(E2E_TESTS) ## Run end to end integration tests

.PHONY: e2e-simple
e2e-simple:
	@echo Running e2e test: simple
	@cd $(CURDIR)/test/e2e/simple && ${GOBIN}/gambol -v run simple.yaml

.PHONY: e2e-hostpath
e2e-hostpath:
	@echo Running e2e test: hostpath
	@cd $(CURDIR)/test/e2e/hostpath && ${GOBIN}/gambol -v run hostpath.yaml

.PHONY: e2e-advanced
e2e-advanced:
	@echo Running e2e test: advanced
	@cd $(CURDIR)/test/e2e/advanced && ${GOBIN}/gambol -v run advanced.yaml

##@ Clean

.PHONY: clean
clean: ## Delete all builds and downloaded dependencies
	@rm -rf bin/
	@rm -rf snap/

FORMATTING_BEGIN_YELLOW = \033[0;33m
FORMATTING_BEGIN_BLUE = \033[36m
FORMATTING_END = \033[0m

.PHONY: help
help:
	@awk 'BEGIN {\
		FS = ":.*##"; \
		printf                "Usage: ${FORMATTING_BEGIN_BLUE}OPTION${FORMATTING_END}=<value> make ${FORMATTING_BEGIN_YELLOW}<target>${FORMATTING_END}\n"\
		} \
		/^[a-zA-Z0-9_-]+:.*?##/ { printf "  ${FORMATTING_BEGIN_BLUE}%-46s${FORMATTING_END} %s\n", $$1, $$2 } \
		/^.?.?##~/              { printf "   %-46s${FORMATTING_BEGIN_YELLOW}%-46s${FORMATTING_END}\n", "", substr($$1, 6) } \
		/^##@/                  { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
