PROJ = gambol
ORG_PATH = github.com/nuccitheboss
REPO_PATH = $(ORG_PATH)/$(PROJ)

export GOBIN=$(CURDIR)/bin

##@ Build

build: bin/gambol ## Build `gambol` binary.

bin/gambol:
	@mkdir -p bin/
	@go install -v $(REPO_PATH)/cmd/gambol

.PHONY: snap
snap: ## Build `gambol` snap package.
	@cp -ra build/snap snap
	@snapcraft -v pack
	@rm -rf snap

##@ Lint

.PHONY: fmt
fmt: ## Format `gambol` source code.
	@go fmt -x ./...
	@golangci-lint run --fix

.PHONY: lint
lint: ## Lint `gambol` source code.
	@golangci-lint run

##@ Clean

.PHONY: clean
clean: ## Delete all builds and downloaded dependencies.
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
