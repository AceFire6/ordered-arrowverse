export GOBIN=$(PWD)/.gobin

GOLANGCI_LINT_VERSION=2.5.0
GOLANGCI_LINT_BIN=$(GOBIN)/golangci-lint

.PHONY: all clean test download-tool-dependencies golangci-lint-version install-linter download-application-dependencies download-dependencies gobin install-pre-commit-hooks install-tools

all: install-tools install-pre-commit-hooks install-linter

clean:
	rm -r "$(GOBIN)"
	go clean

test: download-application-dependencies
	go test ./...

gobin:
	@echo "$(GOBIN)"

golangci-lint-version:
	@echo v$(GOLANGCI_LINT_VERSION)

install-linter:
	@if [ -f "$(GOLANGCI_LINT_BIN)" ] && [ "$$($(GOLANGCI_LINT_BIN) version --short)" = "$(GOLANGCI_LINT_VERSION)" ]; then\
		echo "$(GOLANGCI_LINT_BIN) already installed and at v$(GOLANGCI_LINT_VERSION)";\
	else \
		echo "Installing golangci-lint@$(GOLANGCI_LINT_VERSION)";\
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(GOBIN) "v$(GOLANGCI_LINT_VERSION)";\
	fi

download-tool-dependencies:
	@echo "Downloading tool dependencies"
	cd tools; go mod tidy -v

download-application-dependencies:
	@echo "Downloading application dependencies"
	go mod tidy -v

download-dependencies: download-tool-dependencies download-application-dependencies

install-tools: download-dependencies
	@echo "Installing go tools to $(GOBIN)"
	cd ./tools; go install -v tool

install-pre-commit-hooks:
	@echo "Installing pre-commit hooks"
	@command -v pre-commit > /dev/null 2>&1 || echo "pre-commit missing"
	@pre-commit install

run:
	go run ./cmd/ordered-arrowverse
