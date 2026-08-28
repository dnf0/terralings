.PHONY: all build build-all run test test-race lint fmt-check fmt-solutions verify version-check check clean extension-install extension-build extension-test extension-check extension-package docs-install docs-serve docs-build

all: check build test-race extension-check extension-package

build:
	go build -o bin/terralings ./cmd/terralings
	@if [ "$$(uname -s)" = "Darwin" ]; then codesign -s - -f bin/terralings 2>/dev/null || true; fi

build-all:
	go build -v -o bin/ ./...

run:
	go run ./cmd/terralings

test:
	go test -v ./...

test-race:
	go test -v -race ./...

lint:
	go vet ./...

fmt-check:
	@unformatted=$$(gofmt -l .) || exit 1; \
	if [ -n "$$unformatted" ]; then \
		echo "Unformatted files:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi

fmt-solutions:
	@if command -v tofu >/dev/null 2>&1; then \
		tofu fmt -check -recursive solutions/ || exit 1; \
	elif command -v terraform >/dev/null 2>&1; then \
		terraform fmt -check -recursive solutions/ || exit 1; \
	fi

verify:
	go mod verify

version-check:
	go test -v -run TestVersion ./test

check: verify fmt-check fmt-solutions lint version-check

extension-install:
	cd extensions/vscode && npm install

extension-build:
	cd extensions/vscode && npm run build

extension-test:
	cd extensions/vscode && npm test

extension-check:
	cd extensions/vscode && npm run check-types && npm test

extension-package: extension-build
	mkdir -p dist
	cd extensions/vscode && npx @vscode/vsce package --no-dependencies -o ../../dist/terralings-vscode.vsix

docs-install:
	uv pip install mkdocs-material || pip install mkdocs-material

docs-serve:
	mkdocs serve

docs-build:
	mkdocs build --strict

clean:
	rm -rf bin/ dist/ site/ .terraform/ .cache/ && rm -f ./terralings
