SHELL := /bin/bash
.PHONY: help

help: ## show help
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

build: ## build gojira binary
	@go build -o dist/gojira

install: build
	@sudo cp dist/gojira /usr/local/bin/

tests: ## run tests
	@go test -cover ./...

fixer: ## run static analysis
	@echo "Static analysis..."
	@golangci-lint run --config .golangci.yml --out-format=colored-line-number --concurrency 8