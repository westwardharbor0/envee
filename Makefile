.PHONY: all
all: help

# -- Support functionality --
.PHONY: help dependencies
help: # Show help for each of the Makefile commands.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | \
	while read -r l; do printf "\033[1;32m$$(echo $$l | \
	cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; \
	done

# -- Test --
.PHONY: test
test: # Run local tests.
	go test ./... -cover -race

# -- Local --
.PHONY: lint format
# -- Code Style --
lint: # Run linting to check for issues.
	golangci-lint run

format: # Format code and imports.
	go mod tidy
	gofmt -s -w .
